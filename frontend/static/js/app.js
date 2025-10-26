// 全局变量
let agents = [];
let chatrooms = [];
let currentPage = 'agents';
let editingAgentId = null;
let currentChatroomId = null;
let wsConnection = null;

// 页面初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    
    // 添加头像上传事件监听
    const avatarInput = document.getElementById('agent-avatar');
    if (avatarInput) {
        avatarInput.addEventListener('change', handleAvatarUpload);
    }
});

function initializeApp() {
    setupNavigation();
    setupEventListeners();
    loadData();
    
    // 为测试连接按钮添加事件监听
    const testConnectionBtn = document.getElementById('test-connection');
    if (testConnectionBtn) {
        testConnectionBtn.addEventListener('click', testAPIConnection);
    }
}

// 导航功能
function setupNavigation() {
    const navButtons = document.querySelectorAll('.nav-btn');
    navButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            const page = this.id.replace('nav-', '');
            switchPage(page);
        });
    });
}

function switchPage(page) {
    // 隐藏所有页面
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));
    
    // 显示选中页面
    const pageElement = document.getElementById(page + '-page');
    const navButton = document.getElementById('nav-' + page);
    
    if (pageElement && navButton) {
        pageElement.classList.add('active');
        navButton.classList.add('active');
        currentPage = page;
        
        // 根据页面加载相应数据
        if (page === 'agents') {
            loadAgents();
        } else if (page === 'chatrooms') {
            loadChatrooms();
        }
    }
}

// 事件监听器
function setupEventListeners() {
    // 历史人物管理事件
    const addAgentBtn = document.getElementById('add-agent-btn');
    if (addAgentBtn) {
        addAgentBtn.addEventListener('click', showAgentModal);
    }
    
    const agentForm = document.getElementById('agent-form');
    if (agentForm) {
        agentForm.addEventListener('submit', saveAgent);
    }
    
    // 聊天室事件
    const addChatroomBtn = document.getElementById('add-chatroom-btn');
    if (addChatroomBtn) {
        addChatroomBtn.addEventListener('click', showChatroomModal);
    }
    
    // 模态框关闭事件
    document.querySelectorAll('.close').forEach(closeBtn => {
        closeBtn.addEventListener('click', function() {
            this.closest('.modal').style.display = 'none';
        });
    });
    
    // 点击模态框外部关闭
    document.querySelectorAll('.modal').forEach(modal => {
        modal.addEventListener('click', function(e) {
            if (e.target === this) {
                this.style.display = 'none';
            }
        });
    });
}

// 数据加载
function loadData() {
    loadAgents();
    loadChatrooms();
}

function loadAgents() {
    fetch('/api/agents')
        .then(response => response.json())
        .then(data => {
            agents = data;
            renderAgents();
        })
        .catch(error => {
            console.error('加载历史人物失败:', error);
            showMessage('加载历史人物失败', 'error');
        });
}

function loadChatrooms() {
    fetch('/api/chatrooms')
        .then(response => response.json())
        .then(data => {
            chatrooms = data;
            renderChatrooms();
        })
        .catch(error => {
            console.error('加载聊天室失败:', error);
            showMessage('加载聊天室失败', 'error');
        });
}



// 渲染功能
function renderAgents() {
    const agentsList = document.getElementById('agents-list');
    if (!agentsList) return;
    
    agentsList.innerHTML = '';
    
    agents.forEach(agent => {
        const agentCard = document.createElement('div');
        agentCard.className = 'agent-card';
        agentCard.innerHTML = `
            <div class="agent-avatar-large">
                ${agent.avatar ? `<img src="${agent.avatar}" alt="${agent.name}">` : '<div class="avatar-placeholder">头像</div>'}
            </div>
            <h3>${agent.name}</h3>
            <div class="agent-role">${agent.role}</div>
            <div class="agent-actions">
                <button class="btn btn-secondary" onclick="editAgent('${agent.id}')">编辑</button>
                <button class="btn btn-danger" onclick="deleteAgent('${agent.id}')">删除</button>
                <button class="btn btn-primary" onclick="chatWithAgent('${agent.id}')">对话</button>
            </div>
        `;
        
        // 添加卡片点击事件显示详情
        agentCard.addEventListener('click', function(e) {
            // 如果点击的是按钮，不触发卡片点击事件
            if (e.target.tagName === 'BUTTON') return;
            showAgentDetails(agent.id);
        });
        
        agentsList.appendChild(agentCard);
    });
}

function renderChatrooms() {
    const chatroomsList = document.getElementById('chatrooms-list');
    if (!chatroomsList) return;
    
    chatroomsList.innerHTML = '';
    
    chatrooms.forEach(chatroom => {
        const chatroomCard = document.createElement('div');
        chatroomCard.className = 'chatroom-card';
        
        // 确保主题不为空，如果为空则显示默认文本
        const displayTopic = chatroom.topic && chatroom.topic.trim() !== '' ? chatroom.topic : '未设置主题';
        
        chatroomCard.innerHTML = `
            <h3>${displayTopic}</h3>
            <div class="chatroom-status ${chatroom.status}">${chatroom.status}</div>
            <div class="chatroom-agents-count">${chatroom.agents ? chatroom.agents.length : 0} 个历史人物</div>
            <div class="chatroom-actions">
                <button class="btn btn-primary" onclick="openChatroom('${chatroom.id}')">进入</button>
                <button class="btn btn-danger" onclick="deleteChatroom('${chatroom.id}')">删除</button>
            </div>
        `;
        chatroomsList.appendChild(chatroomCard);
    });
}

// 历史人物管理
function showAgentModal(agentId = null) {
    const modal = document.getElementById('agent-modal');
    const title = document.getElementById('agent-modal-title');
    const form = document.getElementById('agent-form');
    
    // 清空头像预览
    const avatarPreview = document.getElementById('avatar-preview');
    const avatarData = document.getElementById('agent-avatar-data');
    if (avatarPreview) avatarPreview.innerHTML = '';
    if (avatarData) avatarData.value = '';
    
    if (agentId) {
        const agent = agents.find(a => a.id === agentId);
        if (agent) {
            title.textContent = '编辑历史人物';
            document.getElementById('agent-name').value = agent.name;
            document.getElementById('agent-role').value = agent.role;
            document.getElementById('agent-personality').value = agent.personality;
            document.getElementById('agent-coze-link').value = agent.cozeLink || '';
            editingAgentId = agentId;
            
            // 显示现有头像
            if (agent.avatar && avatarPreview) {
                const img = document.createElement('img');
                img.src = agent.avatar;
                img.style.maxWidth = '100px';
                img.style.maxHeight = '100px';
                avatarPreview.appendChild(img);
                if (avatarData) avatarData.value = agent.avatar;
            }
        }
    } else {
        title.textContent = '添加历史人物';
        form.reset();
        editingAgentId = null;
    }
    
    modal.style.display = 'block';
}

// 显示历史人物详情
function showAgentDetails(agentId) {
    const agent = agents.find(a => a.id === agentId);
    if (!agent) return;
    
    const modal = document.getElementById('agent-details-modal');
    const avatar = document.getElementById('details-avatar');
    const name = document.getElementById('details-name');
    const role = document.getElementById('details-role');
    const personality = document.getElementById('details-personality');
    
    if (agent.avatar) {
        avatar.src = agent.avatar;
        avatar.style.display = 'block';
    } else {
        avatar.style.display = 'none';
    }
    
    name.textContent = agent.name;
    role.textContent = agent.role;
    personality.textContent = agent.personality || '暂无性格描述';
    
    modal.style.display = 'block';
}

// 处理头像上传
function handleAvatarUpload(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    // 检查文件类型
    if (!file.type.startsWith('image/')) {
        showMessage('请选择图片文件', 'error');
        return;
    }
    
    // 检查文件大小 (限制为2MB)
    if (file.size > 2 * 1024 * 1024) {
        showMessage('图片大小不能超过2MB', 'error');
        return;
    }
    
    const reader = new FileReader();
    reader.onload = function(e) {
        const avatarPreview = document.getElementById('avatar-preview');
        const avatarData = document.getElementById('agent-avatar-data');
        
        if (avatarPreview) {
            avatarPreview.innerHTML = `<img src="${e.target.result}" style="maxWidth: 100px; maxHeight: 100px;">`;
        }
        
        if (avatarData) {
            avatarData.value = e.target.result;
        }
    };
    
    reader.readAsDataURL(file);
}

function saveAgent(e) {
    e.preventDefault();
    
    console.log('开始保存历史人物...');
    
    const agentData = {
        name: document.getElementById('agent-name').value,
        role: document.getElementById('agent-role').value,
        personality: document.getElementById('agent-personality').value,
        avatar: document.getElementById('agent-avatar-data').value,
        cozeLink: document.getElementById('agent-coze-link').value
    };
    
    console.log('历史人物数据:', agentData);
    
    const url = editingAgentId ? `/api/agents/${editingAgentId}` : '/api/agents';
    const method = editingAgentId ? 'PUT' : 'POST';
    
    console.log('请求URL:', url, '方法:', method);
    
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(agentData)
    })
    .then(response => {
        console.log('保存响应状态:', response.status);
        if (response.ok) {
            return response.json().then(data => {
                console.log('保存成功响应:', data);
                showMessage(editingAgentId ? '历史人物更新成功' : '历史人物添加成功', 'success');
                document.getElementById('agent-modal').style.display = 'none';
                loadAgents();
            });
        } else {
            return response.json().then(data => {
                console.error('保存失败响应:', data);
                showMessage(`保存历史人物失败: ${data.error || '未知错误'}`, 'error');
            }).catch(() => {
                showMessage(`保存历史人物失败: HTTP ${response.status}`, 'error');
            });
        }
    })
    .catch(error => {
        console.error('保存历史人物网络错误:', error);
        showMessage(`保存历史人物失败: ${error.message}`, 'error');
    });
}

function editAgent(agentId) {
    showAgentModal(agentId);
}

function deleteAgent(agentId) {
    if (!agentId) {
        showMessage('无法删除：历史人物ID为空', 'error');
        return;
    }
    
    if (confirm('确定要删除这个历史人物吗？')) {
        console.log('正在删除历史人物，ID:', agentId);
        fetch(`/api/agents/${agentId}`, {
            method: 'DELETE'
        })
        .then(response => {
            console.log('删除响应状态:', response.status);
            if (response.ok) {
                showMessage('历史人物删除成功', 'success');
                loadAgents();
            } else {
                return response.json().then(data => {
                    console.error('删除失败响应:', data);
                    showMessage(`删除历史人物失败: ${data.error || '未知错误'}`, 'error');
                }).catch(() => {
                    showMessage(`删除历史人物失败: HTTP ${response.status}`, 'error');
                });
            }
        })
        .catch(error => {
            console.error('删除历史人物网络错误:', error);
            showMessage(`删除历史人物失败: ${error.message}`, 'error');
        });
    }
}

// 历史圆桌管理
function showChatroomModal() {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <div class="modal-header">
                <h3>创建圆桌</h3>
                <span class="close">&times;</span>
            </div>
            <form id="chatroom-form">
                <div class="form-group">
                    <label for="new-chatroom-topic">讨论主题:</label>
                    <input type="text" id="new-chatroom-topic" placeholder="请输入讨论主题" required>
                </div>
                <div class="form-group">
                    <label>选择历史人物:</label>
                    <div id="agent-selection">
                        ${agents.map(agent => `
                            <label class="agent-checkbox">
                                <input type="checkbox" name="agents" value="${agent.id}">
                                <span>${agent.name} (${agent.role})</span>
                            </label>
                        `).join('')}
                    </div>
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn btn-secondary" onclick="this.closest('.modal').remove()">取消</button>
                    <button type="submit" class="btn btn-primary">创建</button>
                </div>
            </form>
        </div>
    `;
    
    document.body.appendChild(modal);
    modal.style.display = 'block';
    
    // 添加事件监听
    modal.querySelector('.close').addEventListener('click', () => {
        modal.remove();
    });
    
    modal.querySelector('#chatroom-form').addEventListener('submit', function(e) {
        e.preventDefault();
        createChatroom(modal);
    });
    
    // 点击外部关闭
    modal.addEventListener('click', function(e) {
        if (e.target === this) {
            this.remove();
        }
    });
}

function createChatroom(modal) {
    const topicInput = modal.querySelector('#new-chatroom-topic');
    const topic = topicInput ? topicInput.value.trim() : '';
    const selectedAgents = Array.from(modal.querySelectorAll('input[name="agents"]:checked')).map(cb => cb.value);
    
    console.log('创建圆桌 - 主题:', topic, '历史人物:', selectedAgents);
    
    if (!topic) {
        showMessage('请输入讨论主题', 'error');
        return;
    }
    
    if (selectedAgents.length === 0) {
        showMessage('请至少选择一个历史人物', 'error');
        return;
    }
    
    const chatroomData = {
        topic: topic,
        agents: selectedAgents
    };
    
    console.log('发送数据:', chatroomData);
    
    fetch('/api/chatrooms', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(chatroomData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('服务器响应:', data);
        showMessage('圆桌创建成功', 'success');
        modal.remove();
        loadChatrooms();
    })
    .catch(error => {
        console.error('创建圆桌失败:', error);
        showMessage('创建圆桌失败', 'error');
    });
}

function openChatroom(chatroomId) {
    const chatroom = chatrooms.find(c => c.id === chatroomId);
    if (!chatroom) return;
    
    currentChatroomId = chatroomId;
    const modal = document.getElementById('chatroom-modal');
    
    // 确保主题不为空，如果为空则显示默认文本
    const displayTopic = chatroom.topic && chatroom.topic.trim() !== '' ? chatroom.topic : '未设置主题';
    
    document.getElementById('chatroom-title').textContent = displayTopic;
    document.getElementById('chatroom-topic').textContent = displayTopic;
    document.getElementById('chatroom-status').textContent = chatroom.status || 'pending';
    document.getElementById('chatroom-status').className = chatroom.status || 'pending';
    
    // 显示参与历史人物
    const agentsContainer = document.getElementById('chatroom-agents');
    agentsContainer.innerHTML = '';
    
    // 确保chatroom.agents是一个数组
    const chatroomAgents = Array.isArray(chatroom.agents) ? chatroom.agents : [];
    
    chatroomAgents.forEach(agent => {
        // 如果agent是对象（包含完整信息），直接使用
        if (typeof agent === 'object' && agent.name && agent.role) {
            const agentTag = document.createElement('span');
            agentTag.className = 'agent-tag';
            agentTag.textContent = `${agent.name} (${agent.role})`;
            agentsContainer.appendChild(agentTag);
        } else if (typeof agent === 'string') {
            // 如果agent是ID字符串，从agents数组中查找
            const foundAgent = agents.find(a => a.id === agent);
            if (foundAgent) {
                const agentTag = document.createElement('span');
                agentTag.className = 'agent-tag';
                agentTag.textContent = `${foundAgent.name} (${foundAgent.role})`;
                agentsContainer.appendChild(agentTag);
            }
        }
    });
    
    // 清空聊天消息
    const messagesContainer = document.getElementById('chat-messages');
    if (messagesContainer) {
        messagesContainer.innerHTML = '';
    }
    
    modal.style.display = 'block';
    
    // 连接WebSocket
    connectWebSocket(chatroomId);
    
    // 获取聊天室的最新消息
    loadChatroomMessages(chatroomId);
}

function deleteChatroom(chatroomId) {
    if (!chatroomId) {
        showMessage('无法删除：聊天室ID为空', 'error');
        return;
    }
    
    if (confirm('确定要删除这个聊天室吗？')) {
        fetch(`/api/chatrooms/${chatroomId}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (response.ok) {
                showMessage('聊天室删除成功', 'success');
                loadChatrooms();
            } else {
                showMessage('删除聊天室失败', 'error');
            }
        })
        .catch(error => {
            console.error('删除聊天室失败:', error);
            showMessage('删除聊天室失败', 'error');
        });
    }
}

// 加载聊天室消息
function loadChatroomMessages(chatroomId) {
    fetch(`/api/chatrooms/${chatroomId}`)
        .then(response => response.json())
        .then(data => {
            // 显示历史消息
            if (data.messages && Array.isArray(data.messages)) {
                const messagesContainer = document.getElementById('chat-messages');
                if (messagesContainer) {
                    messagesContainer.innerHTML = '';
                    data.messages.forEach(message => {
                        displayMessage(message);
                    });
                }
            }
            
            // 更新聊天室信息
            if (data.topic) {
                const displayTopic = data.topic && data.topic.trim() !== '' ? data.topic : '未设置主题';
                document.getElementById('chatroom-topic').textContent = displayTopic;
                document.getElementById('chatroom-title').textContent = displayTopic;
            }
            if (data.status) {
                const statusElement = document.getElementById('chatroom-status');
                statusElement.textContent = data.status;
                statusElement.className = data.status;
            }
            
            // 更新参与历史人物
            const agentsContainer = document.getElementById('chatroom-agents');
            if (agentsContainer && data.agents && Array.isArray(data.agents)) {
                agentsContainer.innerHTML = '';
                data.agents.forEach(agent => {
                    const agentTag = document.createElement('span');
                    agentTag.className = 'agent-tag';
                    agentTag.textContent = `${agent.name} (${agent.role})`;
                    agentsContainer.appendChild(agentTag);
                });
            }
        })
        .catch(error => {
            console.error('加载聊天室消息失败:', error);
        });
}

// WebSocket连接
function connectWebSocket(chatroomId) {
    if (wsConnection) {
        wsConnection.close();
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/ws/${chatroomId}`;
    
    wsConnection = new WebSocket(wsUrl);
    
    wsConnection.onopen = function() {
        console.log('WebSocket连接已建立');
        document.getElementById('start-discussion').style.display = 'inline-block';
        document.getElementById('stop-discussion').style.display = 'none';
    };
    
    wsConnection.onmessage = function(event) {
        const message = JSON.parse(event.data);
        
        // 处理不同类型的消息
        if (message.type === 'status_update') {
            // 更新聊天室状态
            const statusElement = document.getElementById('chatroom-status');
            if (statusElement) {
                statusElement.textContent = message.status;
                statusElement.className = message.status;
            }
        } else if (message.type === 'room_info') {
            // 更新整个聊天室信息
            if (message.chat_room) {
                const chatroom = message.chat_room;
                document.getElementById('chatroom-topic').textContent = chatroom.topic;
                document.getElementById('chatroom-status').textContent = chatroom.status;
                document.getElementById('chatroom-status').className = chatroom.status;
                
                // 更新参与历史人物
                const agentsContainer = document.getElementById('chatroom-agents');
                if (agentsContainer) {
                    agentsContainer.innerHTML = '';
                    if (chatroom.agents && Array.isArray(chatroom.agents)) {
                        chatroom.agents.forEach(agent => {
                            const agentTag = document.createElement('span');
                            agentTag.className = 'agent-tag';
                            agentTag.textContent = `${agent.name} (${agent.role})`;
                            agentsContainer.appendChild(agentTag);
                        });
                    }
                }
            }
        } else {
            // 显示普通消息
            displayMessage(message);
        }
    };
    
    wsConnection.onclose = function() {
        console.log('WebSocket连接已关闭');
        wsConnection = null;
    };
    
    wsConnection.onerror = function(error) {
        console.error('WebSocket错误:', error);
        showMessage('WebSocket连接失败', 'error');
    };
    
    // 开始讨论按钮事件
    document.getElementById('start-discussion').onclick = function() {
        startDiscussion();
    };
    
    document.getElementById('stop-discussion').onclick = function() {
        stopDiscussion();
    };
}

function startDiscussion() {
    const input = document.getElementById('chat-input').value.trim();
    if (!input) {
        showMessage('请输入讨论话题', 'error');
        return;
    }
    
    // 先调用后端的start接口来激活聊天室
    fetch(`/api/chatrooms/${currentChatroomId}/start`, {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        // 更新聊天室状态
        const statusElement = document.getElementById('chatroom-status');
        if (statusElement) {
            statusElement.textContent = data.status;
            statusElement.className = data.status;
        }
        
        // 然后通过WebSocket发送消息
        if (wsConnection && wsConnection.readyState === WebSocket.OPEN) {
            wsConnection.send(JSON.stringify({
                type: 'start',
                message: input
            }));
            
            document.getElementById('start-discussion').style.display = 'none';
            document.getElementById('stop-discussion').style.display = 'inline-block';
            document.getElementById('chat-input').value = '';
            
            // 显示用户消息
            displayMessage({
                sender: '用户',
                content: input,
                timestamp: new Date().toISOString()
            });
        }
    })
    .catch(error => {
        console.error('启动聊天室失败:', error);
        showMessage('启动聊天室失败', 'error');
    });
}

function stopDiscussion() {
    if (wsConnection && wsConnection.readyState === WebSocket.OPEN) {
        wsConnection.send(JSON.stringify({
            type: 'stop'
        }));
        
        document.getElementById('start-discussion').style.display = 'inline-block';
        document.getElementById('stop-discussion').style.display = 'none';
    }
}

function displayMessage(message) {
    const messagesContainer = document.getElementById('chat-messages');
    if (!messagesContainer) return;
    
    const messageElement = document.createElement('div');
    messageElement.className = 'message';
    
    // 处理不同类型的字段
    const sender = message.sender || message.agent_name || '系统';
    const content = message.content || message.message || '';
    const timestamp = new Date(message.timestamp || Date.now()).toLocaleTimeString();
    
    // 根据消息类型设置样式
    if (message.type === 'agent') {
        messageElement.classList.add('agent-message');
    } else if (message.type === 'user') {
        messageElement.classList.add('user-message');
    } else {
        messageElement.classList.add('system-message');
    }
    
    messageElement.innerHTML = `
        <div class="message-header">
            <span class="message-sender">${sender}</span>
            <span class="message-time">${timestamp}</span>
        </div>
        <div class="message-content">${content}</div>
    `;
    
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// 与历史人物对话
function chatWithAgent(agentId) {
    const agent = agents.find(a => a.id === agentId);
    if (!agent) {
        showMessage('未找到历史人物', 'error');
        return;
    }
    
    // 如果有coze外链，直接跳转到外链
    if (agent.cozeLink) {
        window.open(agent.cozeLink, '_blank');
        return;
    }
    
    // 如果没有外链，显示提示信息
    showMessage('该历史人物暂未配置coze外链', 'warning');
}

// 测试函数
function testEdit() {
    showAgentModal();
    document.getElementById('agent-name').value = '测试历史人物';
    document.getElementById('agent-role').value = '测试助手';
    document.getElementById('agent-personality').value = '这是一个用于测试的历史人物';
}

function testSave() {
    const agentData = {
        name: '测试历史人物',
        role: '测试助手',
        personality: '这是一个用于测试的历史人物'
    };
    
    fetch('/api/agents', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(agentData)
    })
    .then(response => {
        if (response.ok) {
            showMessage('测试历史人物添加成功', 'success');
            loadAgents();
        }
    });
}

// 工具函数
function showMessage(message, type = 'info') {
    // 创建消息提示
    const messageDiv = document.createElement('div');
    messageDiv.className = `message-toast message-${type}`;
    messageDiv.textContent = message;
    
    // 添加样式
    messageDiv.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 20px;
        border-radius: 6px;
        color: white;
        font-weight: 500;
        z-index: 1000;
        max-width: 300px;
        word-wrap: break-word;
        transition: all 0.3s ease;
    `;
    
    // 根据类型设置背景色
    const colors = {
        success: '#48bb78',
        error: '#f56565',
        info: '#4299e1',
        warning: '#ed8936'
    };
    messageDiv.style.backgroundColor = colors[type] || colors.info;
    
    document.body.appendChild(messageDiv);
    
    // 3秒后自动移除
    setTimeout(() => {
        messageDiv.style.opacity = '0';
        setTimeout(() => {
            if (messageDiv.parentNode) {
                messageDiv.parentNode.removeChild(messageDiv);
            }
        }, 300);
    }, 3000);
}

// 测试API连接
function testAPIConnection() {
    const provider = document.getElementById('api-provider').value;
    const apiKey = document.getElementById('api-key').value;
    const baseURL = document.getElementById('base-url').value;
    const model = document.getElementById('model-name').value;
    
    // 验证必填字段
    if (!provider || !apiKey) {
        showMessage('请填写API提供商和密钥', 'error');
        return;
    }
    
    // 显示测试状态
    const testBtn = document.getElementById('test-connection');
    const originalText = testBtn.textContent;
    testBtn.textContent = '测试中...';
    testBtn.disabled = true;
    
    const config = {
        provider: provider,
        apiKey: apiKey,
        baseURL: baseURL,
        model: model
    };
    
    console.log('测试API连接:', config);
    
    // 发送测试请求
    fetch('/api/test-connection', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(config)
    })
    .then(response => response.json())
    .then(data => {
        console.log('测试结果:', data);
        
        if (data.success) {
            showMessage(`连接成功: ${data.message}`, 'success');
            
            // 显示详细信息
            let details = `提供商: ${provider}`;
            if (data.model) {
                details += `\n模型: ${data.model}`;
            }
            if (data.status_code) {
                details += `\n状态码: ${data.status_code}`;
            }
            
            // 3秒后显示详细信息
            setTimeout(() => {
                showMessage(details, 'info');
            }, 2000);
        } else {
            showMessage(`连接失败: ${data.message}`, 'error');
        }
    })
    .catch(error => {
        console.error('测试连接失败:', error);
        showMessage('测试连接失败: ' + error.message, 'error');
    })
    .finally(() => {
        // 恢复按钮状态
        testBtn.textContent = originalText;
        testBtn.disabled = false;
    });
}

// 快速测试API连接功能
function quickTestAPIConnection() {
    // 设置测试数据
    document.getElementById('api-provider').value = 'openai';
    document.getElementById('api-key').value = 'sk-test-key-12345';
    document.getElementById('base-url').value = '';
    document.getElementById('model-name').value = 'gpt-3.5-turbo';
    
    console.log('快速测试：已填充测试数据');
    showMessage('已填充测试数据，请点击"测试连接"按钮', 'info');
}