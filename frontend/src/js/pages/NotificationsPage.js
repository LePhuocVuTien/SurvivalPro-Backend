// Notifications Page Component

function NotificationsPage(app) {
    const notifications = app.filteredData;
    
    // Parse sent status (handle string boolean)
    const notificationsWithStatus = notifications.map(notif => ({
        ...notif,
        sentBool: notif.sent === 'true' || notif.sent === true
    }));
    
    // Parse notification data
    const parseData = (dataStr) => {
        try {
            return JSON.parse(dataStr);
        } catch (e) {
            return {};
        }
    };
    
    // Get notification type icon and color
    const getNotifStyle = (data) => {
        const parsed = parseData(data);
        const styles = {
            'weather': { icon: '🌪️', color: 'bg-red-100 text-red-800', bgColor: 'bg-red-50' },
            'checklist': { icon: '✅', color: 'bg-green-100 text-green-800', bgColor: 'bg-green-50' },
            'guide': { icon: '📚', color: 'bg-purple-100 text-purple-800', bgColor: 'bg-purple-50' },
            'location': { icon: '📍', color: 'bg-blue-100 text-blue-800', bgColor: 'bg-blue-50' },
            'reminder': { icon: '⚠️', color: 'bg-yellow-100 text-yellow-800', bgColor: 'bg-yellow-50' },
            'achievement': { icon: '🎯', color: 'bg-indigo-100 text-indigo-800', bgColor: 'bg-indigo-50' }
        };
        return styles[parsed.type] || { icon: '🔔', color: 'bg-gray-100 text-gray-800', bgColor: 'bg-gray-50' };
    };
    
    return `
        <div class="space-y-6">
            <!-- Stats Cards -->
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Total Notifications</p>
                            <p class="text-2xl font-bold text-gray-800">${app.notifications.length}</p>
                        </div>
                        <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-bell text-orange-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Sent</p>
                            <p class="text-2xl font-bold text-green-600">
                                ${notificationsWithStatus.filter(n => n.sentBool).length}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-check-circle text-green-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Pending</p>
                            <p class="text-2xl font-bold text-yellow-600">
                                ${notificationsWithStatus.filter(n => !n.sentBool).length}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-clock text-yellow-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Users Reached</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${new Set(app.notifications.map(n => n.user_id)).size}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-users text-blue-600 text-xl"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Filters -->
            <div class="bg-white rounded-lg shadow p-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">User</label>
                        <select x-model="filters.notifications.user_id" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500">
                            <option value="">All Users</option>
                            ${app.users.map(user => `
                                <option value="${user.id}">${user.name}</option>
                            `).join('')}
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
                        <select x-model="filters.notifications.sent" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-orange-500">
                            <option value="">All Status</option>
                            <option value="true">Sent</option>
                            <option value="false">Pending</option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Notifications List -->
            <div class="space-y-4">
                ${notificationsWithStatus.length > 0 ? notificationsWithStatus
                    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
                    .map(notif => {
                        const user = app.users.find(u => u.id === notif.user_id);
                        const style = getNotifStyle(notif.data);
                        const data = parseData(notif.data);
                        
                        return `
                            <div class="bg-white rounded-lg shadow hover:shadow-md transition-shadow duration-200">
                                <div class="p-6">
                                    <div class="flex items-start gap-4">
                                        <!-- Icon -->
                                        <div class="w-12 h-12 rounded-full ${style.bgColor} flex items-center justify-center flex-shrink-0">
                                            <span class="text-2xl">${style.icon}</span>
                                        </div>
                                        
                                        <!-- Content -->
                                        <div class="flex-1 min-w-0">
                                            <div class="flex items-start justify-between mb-2">
                                                <h3 class="text-lg font-bold text-gray-800">
                                                    ${notif.title}
                                                </h3>
                                                <div class="flex items-center gap-2 ml-4">
                                                    <span class="px-2 py-1 text-xs font-semibold rounded-full ${notif.sentBool ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}">
                                                        ${notif.sentBool ? 'Sent' : 'Pending'}
                                                    </span>
                                                </div>
                                            </div>
                                            
                                            <p class="text-gray-600 mb-3">
                                                ${notif.body}
                                            </p>
                                            
                                            <!-- Metadata -->
                                            <div class="flex items-center flex-wrap gap-4 text-sm text-gray-500">
                                                <div class="flex items-center">
                                                    <img src="${user?.avatar || 'https://via.placeholder.com/150'}" 
                                                         class="w-6 h-6 rounded-full mr-2"
                                                         onerror="this.src='https://via.placeholder.com/150'">
                                                    <span>${user?.name || 'User ' + notif.user_id}</span>
                                                </div>
                                                <div class="flex items-center">
                                                    <i class="fas fa-clock mr-1"></i>
                                                    ${window.helpers.formatDate(notif.created_at)}
                                                </div>
                                                ${data.type ? `
                                                    <span class="px-2 py-1 text-xs font-medium rounded ${style.color}">
                                                        ${data.type}
                                                    </span>
                                                ` : ''}
                                                ${data.severity ? `
                                                    <span class="px-2 py-1 text-xs font-medium rounded ${
                                                        data.severity === 'high' ? 'bg-red-100 text-red-800' :
                                                        data.severity === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                                                        'bg-blue-100 text-blue-800'
                                                    }">
                                                        ${data.severity} priority
                                                    </span>
                                                ` : ''}
                                            </div>
                                            
                                            <!-- Data Preview -->
                                            ${data && Object.keys(data).length > 1 ? `
                                                <details class="mt-3">
                                                    <summary class="text-xs text-gray-500 cursor-pointer hover:text-gray-700">
                                                        <i class="fas fa-code mr-1"></i>
                                                        View notification data
                                                    </summary>
                                                    <pre class="mt-2 p-3 bg-gray-50 rounded text-xs overflow-x-auto">${JSON.stringify(data, null, 2)}</pre>
                                                </details>
                                            ` : ''}
                                        </div>
                                        
                                        <!-- Actions -->
                                        <div class="flex items-center gap-2">
                                            ${!notif.sentBool ? `
                                                <button @click="sendNotification(${notif.id})" 
                                                        class="px-3 py-2 text-sm bg-green-600 text-white rounded hover:bg-green-700 transition-colors">
                                                    <i class="fas fa-paper-plane mr-1"></i>
                                                    Send
                                                </button>
                                            ` : ''}
                                            <button @click="viewItem('notification', ${JSON.stringify(notif).replace(/"/g, '&quot;')})" 
                                                    class="p-2 text-blue-600 hover:text-blue-900">
                                                <i class="fas fa-eye"></i>
                                            </button>
                                            <button @click="editItem('notification', ${JSON.stringify(notif).replace(/"/g, '&quot;')})" 
                                                    class="p-2 text-indigo-600 hover:text-indigo-900">
                                                <i class="fas fa-edit"></i>
                                            </button>
                                            <button @click="deleteItem('notifications', ${notif.id})" 
                                                    class="p-2 text-red-600 hover:text-red-900">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        `;
                    }).join('') : `
                        <div class="bg-white rounded-lg shadow p-12 text-center">
                            <i class="fas fa-bell-slash text-6xl text-gray-300 mb-4"></i>
                            <p class="text-xl text-gray-500">No notifications found</p>
                            <p class="text-sm text-gray-400 mt-2">Click "Thêm mới" to create a new notification</p>
                        </div>
                    `}
            </div>
        </div>
    `;
}

// Export to window
window.NotificationsPage = NotificationsPage;