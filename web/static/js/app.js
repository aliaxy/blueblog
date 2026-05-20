/* ==========================================================================
   🔵 BlueBlog - Apple HIG Style Single Page App Core Logic
   ========================================================================== */

(function () {
    'use strict';

    // ── 1. 全局应用状态 ──
    const state = {
        token: localStorage.getItem('token') || '',
        username: localStorage.getItem('username') || '',
        userId: localStorage.getItem('user_id') || '',
        currentCommunityId: '',
        currentOrder: 'time', // time or score
        communities: [],
        posts: [],
        votedPosts: JSON.parse(localStorage.getItem('voted_posts') || '{}'), // 本地临时记录用户在此设备上的投票状态
        
        // 分页状态
        currentPage: 1,
        pageSize: 10,
        hasMore: true
    };

    // ── 2. DOM 元素选择器 ──
    const $ = (selector) => document.querySelector(selector);
    const $$ = (selector) => document.querySelectorAll(selector);

    const UI = {
        body: $('body'),
        brandLogo: $('#brand-logo'),
        themeToggle: $('#theme-toggle'),
        btnCreatePostTrigger: $('#btn-create-post-trigger'),
        authNavGroup: $('#auth-nav-group'),
        userNavProfile: $('#user-nav-profile'),
        navUsername: $('#nav-username'),
        navUserAvatar: $('#nav-user-avatar'),
        btnLogout: $('#btn-logout'),
        
        // 帖子流与筛选
        postsContainer: $('#posts-container'),
        sortControl: $('#sort-control'),
        filterIndicator: $('#filter-indicator'),
        currentFilterTag: $('#current-filter-tag'),
        btnClearFilter: $('#btn-clear-filter'),
        communitySidebarList: $('#community-sidebar-list'),
        sidebarWelcomeCard: $('#sidebar-welcome-card'),
        sidebarActionBtn: $('#sidebar-action-btn'),
        btnLoadMore: $('#btn-load-more'), // 分页加载更多按钮
        
        // 登录/注册 Modal
        authModal: $('#auth-modal'),
        btnAuthClose: $('#btn-auth-close'),
        tabLogin: $('#tab-login'),
        tabSignup: $('#tab-signup'),
        formLogin: $('#form-login'),
        formSignup: $('#form-signup'),
        loginUsername: $('#login-username'),
        loginPassword: $('#login-password'),
        signupUsername: $('#signup-username'),
        signupPassword: $('#signup-password'),
        signupRepassword: $('#signup-repassword'),
        loginErrorMsg: $('#login-error-msg'),
        signupErrorMsg: $('#signup-error-msg'),
        loginSpinner: $('#login-spinner'),
        signupSpinner: $('#signup-spinner'),
        btnLoginTrigger: $('#btn-login-trigger'),
        btnSignupTrigger: $('#btn-signup-trigger'),
        
        // 发贴 Modal
        postModal: $('#post-modal'),
        btnPostClose: $('#btn-post-close'),
        formCreatePost: $('#form-create-post'),
        postCommunity: $('#post-community'),
        postTitle: $('#post-title'),
        postContent: $('#post-content'),
        postErrorMsg: $('#post-error-msg'),
        postSubmitSpinner: $('#post-submit-spinner'),
        btnPostCancel: $('#btn-post-cancel'),
        
        // 详情 Modal
        detailModal: $('#detail-modal'),
        btnDetailClose: $('#btn-detail-close'),
        detailCommunityTag: $('#detail-community-tag'),
        detailTitle: $('#detail-title'),
        detailAuthor: $('#detail-author'),
        detailTime: $('#detail-time'),
        detailContent: $('#detail-content'),
        detailVoteUp: $('#detail-vote-up'),
        detailVoteDown: $('#detail-vote-down'),
        detailVoteCount: $('#detail-vote-count')
    };

    // ── 3. API 请求封装 (Fetch Client) ──
    async function apiRequest(path, options = {}) {
        options.headers = options.headers || {};
        if (state.token) {
            options.headers['Authorization'] = `Bearer ${state.token}`;
        }
        options.headers['Content-Type'] = options.headers['Content-Type'] || 'application/json';

        try {
            const response = await fetch(path, options);
            if (response.status === 401 || response.status === 403) {
                // Token 无效或过期，自动登出
                handleLogout();
                showToast('登录状态过期，请重新登录');
                return { code: 1015, msg: '登录已过期' };
            }
            const data = await response.json();
            return data;
        } catch (error) {
            console.error(`API 请求失败 [${path}]:`, error);
            return { code: 500, msg: '无法连接到服务器，请检查网络。' };
        }
    }

    // ── 4. 辅助交互函数 ──
    
    // Apple 模态弹框报错抖动效果
    function triggerModalShake(modalEl) {
        const modalContainer = modalEl.querySelector('.apple-modal');
        if (modalContainer) {
            modalContainer.classList.add('shake');
            setTimeout(() => modalContainer.classList.remove('shake'), 400);
        }
    }

    // 简单的系统通知提示 (Toast)
    function showToast(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = `toast-notice ${type}`;
        toast.style.cssText = `
            position: fixed;
            bottom: 24px;
            left: 50%;
            transform: translateX(-50%) translateY(100px);
            background: rgba(0, 0, 0, 0.85);
            backdrop-filter: blur(10px);
            color: #ffffff;
            padding: 12px 24px;
            border-radius: 12px;
            font-size: 13px;
            font-weight: 500;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            z-index: 9999;
            opacity: 0;
            transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
        `;
        if (type === 'error') {
            toast.style.background = 'rgba(255, 69, 58, 0.95)';
        } else if (type === 'success') {
            toast.style.background = 'rgba(48, 209, 88, 0.95)';
        }
        
        toast.innerText = message;
        document.body.appendChild(toast);
        
        // 触发动画进场
        setTimeout(() => {
            toast.style.opacity = '1';
            toast.style.transform = 'translateX(-50%) translateY(0)';
        }, 50);

        // 退场并销毁
        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateX(-50%) translateY(100px)';
            setTimeout(() => toast.remove(), 400);
        }, 3000);
    }

    // 打开模态窗口
    function openModal(modalEl) {
        modalEl.classList.add('active');
    }

    // 关闭模态窗口
    function closeModal(modalEl) {
        modalEl.classList.remove('active');
        // 清理表单及报错信息
        const form = modalEl.querySelector('form');
        if (form) form.reset();
        const errorMsg = modalEl.querySelector('.form-error');
        if (errorMsg) {
            errorMsg.innerText = '';
            errorMsg.style.display = 'none';
        }
    }

    // 打开写帖子模态框并自动根据当前选中的社区板块进行预选
    function openCreatePostModal() {
        if (state.currentCommunityId) {
            UI.postCommunity.value = state.currentCommunityId;
        } else {
            UI.postCommunity.value = '';
        }
        openModal(UI.postModal);
    }

    // ── 5. 鉴权与用户状态管理 ──
    
    // 初始化登录状态 UI
    function initUserUI() {
        if (state.token && state.username) {
            UI.authNavGroup.style.display = 'none';
            UI.userNavProfile.style.display = 'flex';
            UI.btnCreatePostTrigger.style.display = 'inline-flex';
            
            UI.navUsername.innerText = state.username;
            UI.navUserAvatar.innerText = state.username.substring(0, 1).toUpperCase();
            
            // 侧边栏卡片展示为已登录状态
            UI.sidebarWelcomeCard.innerHTML = `
                <div class="profile-header">
                    <div class="avatar-large" style="background: linear-gradient(135deg, var(--accent-color) 0%, #30d158 100%)">
                        ${state.username.substring(0, 1).toUpperCase()}
                    </div>
                    <div class="welcome-text">
                        <h3 style="font-size: 16px;">你好，${state.username}</h3>
                        <p>欢迎回到极简技术交流社区</p>
                    </div>
                </div>
                <div class="profile-actions">
                    <button class="btn btn-primary btn-block" id="btn-sidebar-create-post">✍️ 发布新的讨论</button>
                </div>
            `;
            
            // 重新绑定发贴事件
            const sidebarCreateBtn = $('#btn-sidebar-create-post');
            if (sidebarCreateBtn) {
                sidebarCreateBtn.addEventListener('click', openCreatePostModal);
            }
        } else {
            UI.authNavGroup.style.display = 'flex';
            UI.userNavProfile.style.display = 'none';
            UI.btnCreatePostTrigger.style.display = 'none';
            
            UI.sidebarWelcomeCard.innerHTML = `
                <div class="profile-header">
                    <div class="avatar-large"></div>
                    <div class="welcome-text">
                        <h3>欢迎来到 BlueBlog</h3>
                        <p>探讨技术、分享日常的极简社区</p>
                    </div>
                </div>
                <div class="profile-actions">
                    <button class="btn btn-primary btn-block" id="sidebar-action-btn">立即登录参与互动</button>
                </div>
            `;
            
            const actionBtn = $('#sidebar-action-btn');
            if (actionBtn) {
                actionBtn.addEventListener('click', () => {
                    switchAuthTab('login');
                    openModal(UI.authModal);
                });
            }
        }
    }

    // 切换登录/注册 Tab
    function switchAuthTab(tabType) {
        if (tabType === 'login') {
            UI.tabLogin.classList.add('active');
            UI.tabSignup.classList.remove('active');
            UI.formLogin.classList.add('active');
            UI.formSignup.classList.remove('active');
        } else {
            UI.tabLogin.classList.remove('active');
            UI.tabSignup.classList.add('active');
            UI.formLogin.classList.remove('active');
            UI.formSignup.classList.add('active');
        }
    }

    // 登录业务处理
    async function handleLogin(e) {
        e.preventDefault();
        UI.loginSpinner.classList.add('active');
        UI.loginErrorMsg.style.display = 'none';

        const username = UI.loginUsername.value.trim();
        const password = UI.loginPassword.value;

        const res = await apiRequest('/api/v1/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });

        UI.loginSpinner.classList.remove('active');

        if (res.code === 1000) {
            showToast('登录成功，欢迎回来！', 'success');
            
            // 存储鉴权数据
            state.token = res.data.token;
            state.username = res.data.user_name || username;
            state.userId = res.data.user_id;
            
            localStorage.setItem('token', state.token);
            localStorage.setItem('username', state.username);
            localStorage.setItem('user_id', state.userId);
            
            initUserUI();
            closeModal(UI.authModal);
            
            // 登录成功后，即时拉取并渲染社区分类
            fetchCommunities();
            // 刷新帖子与投票状态
            fetchPosts();
        } else {
            triggerModalShake(UI.authModal);
            UI.loginErrorMsg.innerText = res.msg || '用户名或密码错误';
            UI.loginErrorMsg.style.display = 'block';
        }
    }

    // 注册业务处理
    async function handleSignup(e) {
        e.preventDefault();
        UI.signupSpinner.classList.add('active');
        UI.signupErrorMsg.style.display = 'none';

        const username = UI.signupUsername.value.trim();
        const password = UI.signupPassword.value;
        const repassword = UI.signupRepassword.value;

        if (password !== repassword) {
            UI.signupSpinner.classList.remove('active');
            triggerModalShake(UI.authModal);
            UI.signupErrorMsg.innerText = '两次输入的密码不一致';
            UI.signupErrorMsg.style.display = 'block';
            return;
        }

        const res = await apiRequest('/api/v1/signup', {
            method: 'POST',
            body: JSON.stringify({
                username,
                password,
                re_password: repassword
            })
        });

        UI.signupSpinner.classList.remove('active');

        if (res.code === 1000) {
            showToast('注册成功！正在为您自动登录...', 'success');
            
            // 注册成功自动尝试登录
            const loginRes = await apiRequest('/api/v1/login', {
                method: 'POST',
                body: JSON.stringify({ username, password })
            });

            if (loginRes.code === 1000) {
                state.token = loginRes.data.token;
                state.username = loginRes.data.user_name || username;
                state.userId = loginRes.data.user_id;
                
                localStorage.setItem('token', state.token);
                localStorage.setItem('username', state.username);
                localStorage.setItem('user_id', state.userId);
                
                initUserUI();
                closeModal(UI.authModal);
                fetchCommunities();
                fetchPosts();
            } else {
                // 登录失败则切换回登录界面让用户手动登录
                switchAuthTab('login');
            }
        } else {
            triggerModalShake(UI.authModal);
            UI.signupErrorMsg.innerText = res.msg || '注册失败，可能用户名已存在';
            UI.signupErrorMsg.style.display = 'block';
        }
    }

    // 退出登录
    function handleLogout() {
        state.token = '';
        state.username = '';
        state.userId = '';
        state.communities = []; // 登出清空社区板块列表缓存
        
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        localStorage.removeItem('user_id');
        
        initUserUI();
        renderCommunitySidebar(); // 清空或更新右侧侧边栏展示
        showToast('您已安全退出登录');
        
        // 刷新列表
        fetchPosts();
    }

    // ── 6. 社区与分类渲染 ──
    async function fetchCommunities() {
        const res = await apiRequest('/api/v1/community', { method: 'GET' });
        if (res.code === 1000 && Array.isArray(res.data)) {
            state.communities = res.data;
            renderCommunitySidebar();
            renderCommunitySelect();
        } else {
            console.error('加载社区板块失败', res);
        }
    }

    // 渲染右侧 macOS 风格侧边栏社区分类
    function renderCommunitySidebar() {
        UI.communitySidebarList.innerHTML = '';
        
        // 追加“全部板块”
        const allItem = document.createElement('div');
        allItem.className = `community-item ${state.currentCommunityId === '' ? 'active' : ''}`;
        allItem.innerHTML = `
            <span>🌍 全部讨论</span>
            <span class="item-arrow">→</span>
        `;
        allItem.addEventListener('click', () => {
            selectCommunity('');
        });
        UI.communitySidebarList.appendChild(allItem);

        // 追加动态社区
        state.communities.forEach(c => {
            const item = document.createElement('div');
            item.className = `community-item ${state.currentCommunityId === c.id.toString() ? 'active' : ''}`;
            item.innerHTML = `
                <span>📁 ${c.name}</span>
                <span class="item-arrow">→</span>
            `;
            item.addEventListener('click', () => {
                selectCommunity(c.id.toString());
            });
            UI.communitySidebarList.appendChild(item);
        });
    }

    // 渲染写帖模态框中的下拉选项
    function renderCommunitySelect() {
        UI.postCommunity.innerHTML = '<option value="" disabled selected>选择目标社区分类</option>';
        state.communities.forEach(c => {
            const option = document.createElement('option');
            option.value = c.id.toString();
            option.innerText = c.name;
            UI.postCommunity.appendChild(option);
        });
    }

    // 激活社区筛选
    function selectCommunity(id) {
        state.currentCommunityId = id;
        
        // 更新侧边栏高亮
        const items = $$('.community-item');
        items.forEach((item, index) => {
            if (id === '' && index === 0) {
                item.classList.add('active');
            } else if (state.communities[index - 1] && state.communities[index - 1].id.toString() === id) {
                item.classList.add('active');
            } else {
                item.classList.remove('active');
            }
        });

        // 重新获取并渲染列表
        renderCommunitySidebar(); // 刷新选中样式
        
        if (id) {
            const selectedCom = state.communities.find(c => c.id.toString() === id);
            UI.currentFilterTag.innerText = selectedCom ? selectedCom.name : '未知';
            UI.filterIndicator.style.display = 'flex';
        } else {
            UI.filterIndicator.style.display = 'none';
        }

        fetchPosts(false); // 切换社区重置为第一页加载
    }

    // ── 7. 帖子流获取与渲染 ──
    async function fetchPosts(append = false) {
        if (!append) {
            state.currentPage = 1;
            renderSkeleton();
            UI.btnLoadMore.style.display = 'none';
        }

        let path = `/api/v1/posts2?page=${state.currentPage}&size=${state.pageSize}&order=${state.currentOrder}`;
        if (state.currentCommunityId) {
            path += `&community_id=${state.currentCommunityId}`;
        }

        const res = await apiRequest(path, { method: 'GET' });
        
        if (res.code === 1000 && Array.isArray(res.data)) {
            if (append) {
                state.posts = state.posts.concat(res.data);
            } else {
                state.posts = res.data;
            }

            // 判断是否还有更多页数据
            if (res.data.length < state.pageSize) {
                state.hasMore = false;
                UI.btnLoadMore.style.display = 'none';
            } else {
                state.hasMore = true;
                UI.btnLoadMore.style.display = 'inline-block';
            }
            
            renderPostsFeed();
        } else {
            if (!append) {
                UI.postsContainer.innerHTML = `
                    <div style="text-align: center; padding: 40px 20px; color: var(--text-secondary);">
                        <p style="font-size: 14px;">📭 暂无相关讨论贴，快去发布第一篇吧！</p>
                    </div>
                `;
            }
            UI.btnLoadMore.style.display = 'none';
        }
    }

    // 骨架屏渲染
    function renderSkeleton() {
        UI.postsContainer.innerHTML = Array(3).fill(`
            <div class="skeleton-card">
                <div class="skeleton-vote"></div>
                <div class="skeleton-content">
                    <div class="skeleton-line skeleton-title"></div>
                    <div class="skeleton-line skeleton-text"></div>
                    <div class="skeleton-line skeleton-meta"></div>
                </div>
            </div>
        `).join('');
    }

    // 动态渲染帖子列表
    function renderPostsFeed() {
        if (state.posts.length === 0) {
            UI.postsContainer.innerHTML = `
                <div style="text-align: center; padding: 40px 20px; color: var(--text-secondary);">
                    <p style="font-size: 14px;">📭 该板块暂无讨论帖子...</p>
                </div>
            `;
            return;
        }

        UI.postsContainer.innerHTML = '';
        state.posts.forEach(post => {
            const card = document.createElement('article');
            card.className = 'post-card';
            
            // 解析时间
            const timeStr = formatTime(post.create_time);
            
            // 获取用户当前对此帖子的本地投票历史记录
            const localVote = state.votedPosts[post.id] || 0;
            const upActive = localVote === 1 ? 'active' : '';
            const downActive = localVote === -1 ? 'active' : '';
            
            const comName = post.community ? post.community.name : '公共论坛';
            const excerpt = post.content.length > 120 ? post.content.substring(0, 120) + '...' : post.content;
            
            card.innerHTML = `
                <div class="post-vote-bar" data-id="${post.id}">
                    <button class="vote-btn up ${upActive}" title="赞同">▲</button>
                    <span class="vote-count">${post.vote_num || 0}</span>
                    <button class="vote-btn down ${downActive}" title="反对">▼</button>
                </div>
                <div class="post-content-area">
                    <div>
                        <div class="post-card-header">
                            <span class="tag">${comName}</span>
                            <span style="font-size: 11px; color: var(--text-secondary);">发布于 ${timeStr}</span>
                        </div>
                        <h2 class="post-card-title">${escapeHTML(post.title)}</h2>
                        <p class="post-card-excerpt">${escapeHTML(excerpt)}</p>
                    </div>
                    <div class="post-card-meta">
                        <span>✍️ 作者: <span class="author-badge">${escapeHTML(post.author_name)}</span></span>
                        <span>👍 ${post.vote_num || 0} 点赞</span>
                    </div>
                </div>
            `;
            
            // 绑定点击内容区域打开深读模态框
            card.querySelector('.post-content-area').addEventListener('click', () => {
                showPostDetail(post);
            });
            
            // 绑定点赞点踩事件
            const voteBar = card.querySelector('.post-vote-bar');
            voteBar.querySelector('.vote-btn.up').addEventListener('click', (e) => {
                e.stopPropagation();
                handleVoteClick(post.id, 1, voteBar);
            });
            voteBar.querySelector('.vote-btn.down').addEventListener('click', (e) => {
                e.stopPropagation();
                handleVoteClick(post.id, -1, voteBar);
            });

            UI.postsContainer.appendChild(card);
        });
    }

    // ── 8. 投票交互处理 ──
    async function handleVoteClick(postId, direction, voteBarEl) {
        if (!state.token) {
            showToast('请先登录系统后再进行投票', 'error');
            switchAuthTab('login');
            openModal(UI.authModal);
            return;
        }

        const currentLocalDirection = state.votedPosts[postId] || 0;
        let finalDirection = direction;
        
        // 如果点击已激活的按钮，则表示取消投票
        if (currentLocalDirection === direction) {
            finalDirection = 0;
        }

        // 调用后端投票 API
        const res = await apiRequest('/api/v1/vote', {
            method: 'POST',
            body: JSON.stringify({
                post_id: postId.toString(),
                direction: finalDirection
            })
        });

        if (res.code === 1000) {
            // 本地状态更新
            state.votedPosts[postId] = finalDirection;
            localStorage.setItem('voted_posts', JSON.stringify(state.votedPosts));
            
            // 动态更新卡片上的投票 UI (免去重新刷新整页)
            const upBtn = voteBarEl.querySelector('.vote-btn.up');
            const downBtn = voteBarEl.querySelector('.vote-btn.down');
            const countEl = voteBarEl.querySelector('.vote-count');
            
            let originalCount = parseInt(countEl.innerText) || 0;
            
            // 减去之前的贡献分，加上现在的贡献分以实现平滑动效
            let delta = finalDirection - currentLocalDirection;
            // 注意：后端的实际 ZSet 分数增加可能存在复杂规则，但这里前端展示进行实时的 +1/-1/0 纯赞同数值渲染
            let displayDelta = 0;
            if (currentLocalDirection === 1) displayDelta -= 1;
            if (currentLocalDirection === -1) displayDelta += 1;
            if (finalDirection === 1) displayDelta += 1;
            if (finalDirection === -1) displayDelta -= 1;

            const newCount = originalCount + displayDelta;
            countEl.innerText = newCount;

            // 切换激活状态
            upBtn.classList.remove('active');
            downBtn.classList.remove('active');
            if (finalDirection === 1) {
                upBtn.classList.add('active');
                showToast('已赞同该贴');
            } else if (finalDirection === -1) {
                downBtn.classList.add('active');
                showToast('已反对该贴');
            } else {
                showToast('已取消投票');
            }
            
            // 同步更新已加载内存中相应帖子对象的赞同数
            const postObj = state.posts.find(p => p.id === postId);
            if (postObj) {
                postObj.vote_num = newCount;
            }
        } else {
            showToast(res.msg || '投票提交失败', 'error');
        }
    }

    // ── 9. 创建新帖子业务 ──
    async function handleCreatePost(e) {
        e.preventDefault();
        UI.postSubmitSpinner.classList.add('active');
        UI.postErrorMsg.style.display = 'none';

        const communityId = UI.postCommunity.value;
        const title = UI.postTitle.value.trim();
        const content = UI.postContent.value.trim();

        if (!communityId) {
            UI.postSubmitSpinner.classList.remove('active');
            triggerModalShake(UI.postModal);
            UI.postErrorMsg.innerText = '请选择一个社区板块发布';
            UI.postErrorMsg.style.display = 'block';
            return;
        }

        const res = await apiRequest('/api/v1/post', {
            method: 'POST',
            body: JSON.stringify({
                title,
                content,
                community_id: communityId // 发送 string 契合 Snowflake ID 映射
            })
        });

        UI.postSubmitSpinner.classList.remove('active');

        if (res.code === 1000) {
            showToast('帖子创建并发布成功！', 'success');
            closeModal(UI.postModal);
            // 自动刷新主贴流，置为第一页并覆写
            fetchPosts(false);
        } else {
            triggerModalShake(UI.postModal);
            UI.postErrorMsg.innerText = res.msg || '发布帖子失败，请重试';
            UI.postErrorMsg.style.display = 'block';
        }
    }

    // ── 10. 帖子深读详情展示 ──
    async function showPostDetail(post) {
        // A. 瞬时极速交互：利用已有缓存数据做首屏无缝加载渲染，带来零卡顿秒开体验
        renderDetailUI(post);
        openModal(UI.detailModal);

        // B. 动态升级：如果用户处于已登录状态，后台默默拉取接口细节最新状态（如最新点赞数、最全大文本等）
        if (state.token) {
            const res = await apiRequest(`/api/v1/post/${post.id}`, { method: 'GET' });
            if (res.code === 1000 && res.data) {
                // 如果发现有最新正文数据，进行无感替换渲染
                renderDetailUI(res.data);
            }
        }
    }

    // 专门负责大文章详情页渲染的界面渲染函数
    function renderDetailUI(postDetail) {
        const comName = postDetail.community ? postDetail.community.name : '公共论坛';
        
        UI.detailCommunityTag.innerText = comName;
        UI.detailTitle.innerText = postDetail.title;
        UI.detailAuthor.innerText = postDetail.author_name || '匿名作者';
        UI.detailTime.innerText = `发布于 ${formatTime(postDetail.create_time)}`;
        
        //  Apple 高端排版：使用微型高兼容 Markdown 编译器进行富文本转换！
        UI.detailContent.innerHTML = compileMarkdown(postDetail.content);
        
        // 更新投票数量与高亮激活态
        UI.detailVoteCount.innerText = postDetail.vote_num || 0;
        
        const localVote = state.votedPosts[postDetail.id] || 0;
        UI.detailVoteUp.className = `vote-action-btn up ${localVote === 1 ? 'active' : ''}`;
        UI.detailVoteDown.className = `vote-action-btn down ${localVote === -1 ? 'active' : ''}`;

        // 重新绑定详情界面的点赞点踩事件
        UI.detailVoteUp.onclick = () => {
            const detailVoteBarFake = {
                querySelector: (sel) => {
                    if (sel.includes('up')) return UI.detailVoteUp;
                    if (sel.includes('down')) return UI.detailVoteDown;
                    if (sel.includes('count')) return UI.detailVoteCount;
                }
            };
            handleVoteClick(postDetail.id, 1, detailVoteBarFake);
        };

        UI.detailVoteDown.onclick = () => {
            const detailVoteBarFake = {
                querySelector: (sel) => {
                    if (sel.includes('up')) return UI.detailVoteUp;
                    if (sel.includes('down')) return UI.detailVoteDown;
                    if (sel.includes('count')) return UI.detailVoteCount;
                }
            };
            handleVoteClick(postDetail.id, -1, detailVoteBarFake);
        };
    }

    // ── 11. 工具函数 ──
    
    // HTML 转义防 XSS 注入
    function escapeHTML(str) {
        if (!str) return '';
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }

    //  极简优雅、高安全的内置微型 Markdown 编译器 (Satisfies technical blog requirements)
    function compileMarkdown(text) {
        if (!text) return '';
        // 1. 先进行 XSS 转义防护
        let safeHTML = escapeHTML(text);

        // 2. 编译代码块 ```code```
        safeHTML = safeHTML.replace(/```([\s\S]*?)```/g, '<pre style="background: var(--input-bg); padding: 14px; border-radius: 10px; font-family: monospace; overflow-x: auto; margin: 12px 0; border: 1px solid var(--border-color); line-height: 1.4;"><code>$1</code></pre>');

        // 3. 编译行内代码 `code`
        safeHTML = safeHTML.replace(/`([^`]+)`/g, '<code style="background: var(--input-bg); padding: 2px 6px; border-radius: 5px; font-family: monospace; font-size: 0.9em; border: 1px solid var(--border-color); color: var(--accent-color);">$1</code>');

        // 4. 编译粗体 **bold** 与 斜体 *italic*
        safeHTML = safeHTML.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
        safeHTML = safeHTML.replace(/\*([^*]+)\*/g, '<em>$1</em>');

        // 5. 编译链接 [title](url)
        safeHTML = safeHTML.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" style="color: var(--accent-color); text-decoration: underline;">$1</a>');

        // 6. 编译段落/换行
        safeHTML = safeHTML.replace(/\n/g, '<br>');

        return safeHTML;
    }

    // 美化时间呈现
    function formatTime(isoString) {
        if (!isoString) return '刚刚';
        const date = new Date(isoString);
        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMins / 60);
        
        if (diffMins < 1) return '刚刚';
        if (diffMins < 60) return `${diffMins} 分钟前`;
        if (diffHours < 24) return `${diffHours} 小时前`;
        
        // 格式化输出为 2026-05-20
        const y = date.getFullYear();
        const m = String(date.getMonth() + 1).padStart(2, '0');
        const d = String(date.getDate()).padStart(2, '0');
        return `${y}-${m}-${d}`;
    }

    // ── 12. 事件绑定与初始化 ──
    function initEvents() {
        // 主题切换
        UI.themeToggle.addEventListener('click', toggleTheme);
        
        // Logo 点击回到全部帖子
        UI.brandLogo.addEventListener('click', () => selectCommunity(''));
        
        // 分页“加载更多”按钮事件
        UI.btnLoadMore.addEventListener('click', () => {
            state.currentPage++;
            fetchPosts(true); // 传入 true 表示追加渲染帖子列表
        });
        
        // 授权模态框事件
        UI.btnLoginTrigger.addEventListener('click', () => {
            switchAuthTab('login');
            openModal(UI.authModal);
        });
        UI.btnSignupTrigger.addEventListener('click', () => {
            switchAuthTab('signup');
            openModal(UI.authModal);
        });
        UI.btnAuthClose.addEventListener('click', () => closeModal(UI.authModal));
        UI.tabLogin.addEventListener('click', () => switchAuthTab('login'));
        UI.tabSignup.addEventListener('click', () => switchAuthTab('signup'));
        
        UI.formLogin.addEventListener('submit', handleLogin);
        UI.formSignup.addEventListener('submit', handleSignup);
        UI.btnLogout.addEventListener('click', handleLogout);

        // 写发贴模态框事件
        UI.btnCreatePostTrigger.addEventListener('click', openCreatePostModal);
        UI.btnPostClose.addEventListener('click', () => closeModal(UI.postModal));
        UI.btnPostCancel.addEventListener('click', () => closeModal(UI.postModal));
        UI.formCreatePost.addEventListener('submit', handleCreatePost);

        // 详情模态框关闭
        UI.btnDetailClose.addEventListener('click', () => closeModal(UI.detailModal));
        
        // 模态框背景层点击关闭
        [UI.authModal, UI.postModal, UI.detailModal].forEach(modal => {
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    closeModal(modal);
                }
            });
        });

        // 排序 Tab 扁平分段胶囊控件事件
        const sortRadios = $$('input[name="sort"]');
        sortRadios.forEach(radio => {
            radio.addEventListener('change', (e) => {
                state.currentOrder = e.target.value;
                
                // 平移动画滑块宽度适应
                const slider = $('.segment-slider');
                if (state.currentOrder === 'time') {
                    slider.style.transform = 'translateX(0)';
                } else {
                    slider.style.transform = 'translateX(88px)';
                }
                
                fetchPosts(false); // 重置列表并重新查询
            });
        });

        // 清理筛选条件标签
        UI.btnClearFilter.addEventListener('click', () => {
            selectCommunity('');
        });
    }

    // 极简主题控制逻辑
    function initTheme() {
        const savedTheme = localStorage.getItem('theme') || 'light';
        applyTheme(savedTheme);
    }

    function toggleTheme() {
        const isDark = UI.body.classList.contains('dark-theme');
        const nextTheme = isDark ? 'light' : 'dark';
        applyTheme(nextTheme);
    }

    function applyTheme(theme) {
        if (theme === 'dark') {
            UI.body.classList.remove('light-theme');
            UI.body.classList.add('dark-theme');
            UI.themeToggle.innerText = '☀️';
            localStorage.setItem('theme', 'dark');
        } else {
            UI.body.classList.remove('dark-theme');
            UI.body.classList.add('light-theme');
            UI.themeToggle.innerText = '🌓';
            localStorage.setItem('theme', 'light');
        }
    }

    // ── 13. 系统引导入口 ──
    function init() {
        initTheme();
        initUserUI();
        initEvents();
        
        // 数据首屏装载：仅在已登录且存在 Token 时获取社区分类，防止 401 报错
        if (state.token) {
            fetchCommunities();
        } else {
            UI.communitySidebarList.innerHTML = '<div style="font-size: 12px; color: var(--text-secondary); text-align: center; padding: 16px 0; font-weight: 500;">🔒 登录后解锁社区讨论板块</div>';
        }
        fetchPosts(false);
    }

    // DOM 构建完成开始执行
    document.addEventListener('DOMContentLoaded', init);

})();
