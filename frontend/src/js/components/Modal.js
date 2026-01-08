// Modal component

function ModalComponent(app) {
    // Form templates for different modal types
    const formTemplates = {
        user: (data, users) => `
            <form @submit.prevent="saveItem()" class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Tên</label>
                        <input type="text" x-model="modal.data.name" required
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <input type="email" x-model="modal.data.email" required
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                    </div>
                </div>
                
                ${!data.id ? `
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Mật khẩu</label>
                        <input type="password" x-model="modal.data.password" required
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                    </div>
                ` : ''}
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Avatar URL</label>
                    <input type="text" x-model="modal.data.avatar"
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <div class="mt-2 flex items-center">
                        <div class="w-12 h-12 rounded-full overflow-hidden bg-gray-200 mr-3">
                            <img :src="modal.data.avatar || 'https://via.placeholder.com/150'" 
                                 class="w-full h-full object-cover"
                                 @error="modal.data.avatar = 'https://via.placeholder.com/150'">
                        </div>
                        <span class="text-sm text-gray-500">Preview</span>
                    </div>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Push Token</label>
                    <input type="text" x-model="modal.data.push_token"
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
            </form>
        `,

        checklist: (data, users) => `
            <form @submit.prevent="saveItem()" class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">User</label>
                    <select x-model="modal.data.user_id" required
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option value="">Chọn user</option>
                        ${users.map(user => `
                            <option value="${user.id}" ${data.user_id == user.id ? 'selected' : ''}>${user.name}</option>
                        `).join('')}
                    </select>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Tiêu đề</label>
                    <input type="text" x-model="modal.data.title" required
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Danh mục</label>
                    <select x-model="modal.data.category" required
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option value="supplies" ${data.category === 'supplies' ? 'selected' : ''}>Supplies</option>
                        <option value="documents" ${data.category === 'documents' ? 'selected' : ''}>Documents</option>
                        <option value="emergency" ${data.category === 'emergency' ? 'selected' : ''}>Emergency</option>
                        <option value="food" ${data.category === 'food' ? 'selected' : ''}>Food</option>
                        <option value="water" ${data.category === 'water' ? 'selected' : ''}>Water</option>
                        <option value="shelter" ${data.category === 'shelter' ? 'selected' : ''}>Shelter</option>
                    </select>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Mô tả</label>
                    <textarea x-model="modal.data.description" rows="3"
                              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">${data.description || ''}</textarea>
                </div>
                
                <div class="flex items-center">
                    <input type="checkbox" x-model="modal.data.is_checked" ${data.is_checked ? 'checked' : ''}
                           class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500">
                    <label class="ml-2 text-sm font-medium text-gray-700">Đã hoàn thành</label>
                </div>
            </form>
        `,

        guide: (data, users) => `
            <form @submit.prevent="saveItem()" class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">User</label>
                    <select x-model="modal.data.user_id" required
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option value="">Chọn user</option>
                        ${users.map(user => `
                            <option value="${user.id}" ${data.user_id == user.id ? 'selected' : ''}>${user.name}</option>
                        `).join('')}
                    </select>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Tiêu đề</label>
                    <input type="text" x-model="modal.data.title" required
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Danh mục</label>
                        <select x-model="modal.data.category" required
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="first-aid" ${data.category === 'first-aid' ? 'selected' : ''}>First Aid</option>
                            <option value="shelter" ${data.category === 'shelter' ? 'selected' : ''}>Shelter</option>
                            <option value="food" ${data.category === 'food' ? 'selected' : ''}>Food & Water</option>
                            <option value="navigation" ${data.category === 'navigation' ? 'selected' : ''}>Navigation</option>
                            <option value="fire" ${data.category === 'fire' ? 'selected' : ''}>Fire Making</option>
                            <option value="signaling" ${data.category === 'signaling' ? 'selected' : ''}>Signaling</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Độ khó</label>
                        <select x-model="modal.data.difficulty" required
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="easy" ${data.difficulty === 'easy' ? 'selected' : ''}>Easy</option>
                            <option value="medium" ${data.difficulty === 'medium' ? 'selected' : ''}>Medium</option>
                            <option value="hard" ${data.difficulty === 'hard' ? 'selected' : ''}>Hard</option>
                        </select>
                    </div>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Icon</label>
                    <input type="text" x-model="modal.data.icon" required
                           placeholder="🔪, 🏕️, 🔥, v.v."
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Nội dung</label>
                    <textarea x-model="modal.data.content" rows="4" required
                              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">${data.content || ''}</textarea>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Hình ảnh URL</label>
                    <input type="text" x-model="modal.data.image_url"
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                
                <div class="flex items-center">
                    <input type="checkbox" x-model="modal.data.is_read" ${data.is_read ? 'checked' : ''}
                           class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500">
                    <label class="ml-2 text-sm font-medium text-gray-700">Đã đọc</label>
                </div>
            </form>
        `,

        notification: (data, users) => `
            <form @submit.prevent="saveItem()" class="space-y-4">
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">User</label>
                    <select x-model="modal.data.user_id" required
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option value="">Chọn user</option>
                        ${users.map(user => `
                            <option value="${user.id}" ${data.user_id == user.id ? 'selected' : ''}>${user.name}</option>
                        `).join('')}
                    </select>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Tiêu đề</label>
                    <input type="text" x-model="modal.data.title" required
                           class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Nội dung</label>
                    <textarea x-model="modal.data.body" rows="3" required
                              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">${data.body || ''}</textarea>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Loại</label>
                        <select x-model="modal.data.type" required
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="push" ${data.type === 'push' ? 'selected' : ''}>Push</option>
                            <option value="in_app" ${data.type === 'in_app' ? 'selected' : ''}>In App</option>
                            <option value="both" ${data.type === 'both' ? 'selected' : ''}>Both</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Data (JSON)</label>
                        <input type="text" x-model="modal.data.data"
                               placeholder='{"key": "value"}'
                               class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                    </div>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex items-center">
                        <input type="checkbox" x-model="modal.data.is_read" ${data.is_read ? 'checked' : ''}
                               class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500">
                        <label class="ml-2 text-sm font-medium text-gray-700">Đã đọc</label>
                    </div>
                    <div class="flex items-center">
                        <input type="checkbox" x-model="modal.data.sent" ${data.sent ? 'checked' : ''}
                               class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500">
                        <label class="ml-2 text-sm font-medium text-gray-700">Đã gửi</label>
                    </div>
                </div>
            </form>
        `
    };

    // View templates for different modal types
    const viewTemplates = {
        user: (data) => `
            <div>
                <div class="flex items-center mb-4">
                    <img src="${data.avatar || 'https://via.placeholder.com/150'}" 
                         class="w-20 h-20 rounded-full object-cover border-2 border-gray-300 mr-4"
                         onerror="this.src='https://via.placeholder.com/150'">
                    <div>
                        <h4 class="text-xl font-bold text-gray-800">${data.name}</h4>
                        <p class="text-gray-600">${data.email}</p>
                    </div>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p class="text-sm text-gray-500">Push Token</p>
                        <p class="font-medium">${data.push_token || 'Không có'}</p>
                    </div>
                    <div>
                        <p class="text-sm text-gray-500">Ngày tạo</p>
                        <p class="font-medium">${window.helpers.formatDate(data.created_at)}</p>
                    </div>
                </div>
            </div>
        `,

        checklist: (data) => `
            <div class="space-y-4">
                <div>
                    <h4 class="text-xl font-bold text-gray-800 mb-2">${data.title}</h4>
                    <div class="flex items-center mb-4">
                        <span class="px-3 py-1 text-xs font-medium rounded-full mr-2 ${window.helpers.getBadgeClassByCategory(data.category)}">
                            ${data.category}
                        </span>
                        <span class="px-2 py-1 text-xs font-medium rounded-full ${data.is_checked ? 'badge-green' : 'badge-yellow'}">
                            ${data.is_checked ? 'Đã hoàn thành' : 'Chưa hoàn thành'}
                        </span>
                    </div>
                </div>
                <div>
                    <p class="text-sm text-gray-500 mb-1">Mô tả</p>
                    <p class="text-gray-800">${data.description || 'Không có mô tả'}</p>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p class="text-sm text-gray-500">Người tạo</p>
                        <p class="font-medium">${window.helpers.getUserName(app.users, data.user_id)}</p>
                    </div>
                    <div>
                        <p class="text-sm text-gray-500">Ngày tạo</p>
                        <p class="font-medium">${window.helpers.formatDate(data.created_at)}</p>
                    </div>
                </div>
            </div>
        `
        // ... other view templates
    };

    return {
        getModalContent() {
            if (app.modal.mode === 'view') {
                const template = viewTemplates[app.modal.type];
                return template ? template(app.modal.data) : '<p>No view template available</p>';
            } else {
                const template = formTemplates[app.modal.type];
                return template ? template(app.modal.data, app.users) : '<p>No form template available</p>';
            }
        }
    };
}

// Export modal component
window.ModalComponent = ModalComponent;