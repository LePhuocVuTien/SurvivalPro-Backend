// Formatter functions for data display

function formatUserRow(user, app) {
    return `
        <tr class="table-row-hover">
            <td class="px-6 py-4 text-sm text-gray-900">${user.id}</td>
            <td class="px-6 py-4">
                <img src="${user.avatar || 'https://via.placeholder.com/40'}" 
                     class="w-10 h-10 rounded-full object-cover border-2 border-gray-200"
                     onerror="this.src='https://via.placeholder.com/40'">
            </td>
            <td class="px-6 py-4 text-sm font-medium text-gray-900">${user.name}</td>
            <td class="px-6 py-4 text-sm text-gray-600">${user.email}</td>
            <td class="px-6 py-4 text-sm text-gray-600">${window.helpers.formatDate(user.created_at)}</td>
            <td class="px-6 py-4 text-sm space-x-3">
                <button onclick="app.viewItem('user', ${JSON.stringify(user).replace(/"/g, '&quot;')})" 
                        class="text-blue-600 hover:text-blue-800 transition-colors" 
                        title="View">
                    <i class="fas fa-eye"></i>
                </button>
                <button onclick="app.editItem('user', ${JSON.stringify(user).replace(/"/g, '&quot;')})" 
                        class="text-green-600 hover:text-green-800 transition-colors" 
                        title="Edit">
                    <i class="fas fa-edit"></i>
                </button>
                <button onclick="app.deleteItem('users', ${user.id})" 
                        class="text-red-600 hover:text-red-800 transition-colors" 
                        title="Delete">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        </tr>
    `;
}

function formatChecklistRow(item, app) {
    const isChecked = item.is_checked;
    const badgeClass = window.helpers.getBadgeClassByCategory(item.category);
    
    return `
        <tr class="table-row-hover">
            <td class="px-6 py-4">
                <input type="checkbox" 
                       ${isChecked ? 'checked' : ''}
                       onclick="app.toggleChecked(${item.id})"
                       class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500 cursor-pointer">
            </td>
            <td class="px-6 py-4 text-sm font-medium text-gray-900 ${isChecked ? 'line-through text-gray-400' : ''}">${item.title}</td>
            <td class="px-6 py-4">
                <span class="px-3 py-1 text-xs font-medium rounded-full ${badgeClass}">
                    ${item.category}
                </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-600 line-clamp-2 max-w-xs">${item.description || 'Không có mô tả'}</td>
            <td class="px-6 py-4 text-sm text-gray-600">${window.helpers.getUserName(app.users, item.user_id)}</td>
            <td class="px-6 py-4 text-sm space-x-3">
                <button onclick="app.viewItem('checklist', ${JSON.stringify(item).replace(/"/g, '&quot;')})" 
                        class="text-blue-600 hover:text-blue-800 transition-colors">
                        <i class="fas fa-eye"></i>
                </button>
                <button onclick="app.editItem('checklist', ${JSON.stringify(item).replace(/"/g, '&quot;')})" 
                        class="text-green-600 hover:text-green-800 transition-colors">
                        <i class="fas fa-edit"></i>
                </button>
                <button onclick="app.deleteItem('checklist', ${item.id})" 
                        class="text-red-600 hover:text-red-800 transition-colors">
                        <i class="fas fa-trash"></i>
                </button>
            </td>
        </tr>
    `;
}

function formatLocationRow(location, app) {
    return `
        <tr class="hover:bg-blue-50 cursor-pointer transition-colors" 
            onclick="app.focusLocation(${JSON.stringify(location).replace(/"/g, '&quot;')})">
            <td class="px-4 py-3 text-sm font-medium text-gray-900">${window.helpers.getUserName(app.users, location.user_id)}</td>
            <td class="px-4 py-3 text-sm text-gray-600 font-mono">${location.lat.toFixed(6)}</td>
            <td class="px-4 py-3 text-sm text-gray-600 font-mono">${location.lon.toFixed(6)}</td>
            <td class="px-4 py-3 text-sm text-gray-600">${window.helpers.formatDate(location.created_at)}</td>
        </tr>
    `;
}

function formatGuideCard(guide, app) {
    const difficultyBadge = window.helpers.getDifficultyBadgeClass(guide.difficulty);
    const categoryBadge = window.helpers.getBadgeClassByCategory(guide.category);
    
    return `
        <div class="bg-white rounded-lg shadow overflow-hidden card-hover">
            <div class="relative">
                <img src="${guide.image_url || 'https://via.placeholder.com/800x400'}" 
                     class="w-full h-48 object-cover"
                     onerror="this.src='https://via.placeholder.com/800x400'">
                <div class="absolute top-2 right-2 bg-white rounded-full px-2 py-1 text-xs font-medium shadow">
                    ${guide.icon}
                </div>
            </div>
            <div class="p-4">
                <div class="flex items-center gap-2 mb-2 flex-wrap">
                    <span class="px-2 py-1 text-xs font-medium rounded-full ${difficultyBadge}">
                        ${guide.difficulty.toUpperCase()}
                    </span>
                    <span class="px-2 py-1 text-xs font-medium rounded-full ${categoryBadge} capitalize">
                        ${guide.category}
                    </span>
                    ${guide.is_read ? '<span class="px-2 py-1 text-xs font-medium rounded-full badge-green">✓ Đã đọc</span>' : ''}
                </div>
                <h3 class="font-semibold text-gray-800 mb-2 line-clamp-2">${guide.title}</h3>
                <p class="text-sm text-gray-600 mb-3 capitalize">${guide.category}</p>
                <div class="flex items-center justify-between text-xs text-gray-500 mb-3">
                    <span><i class="fas fa-eye mr-1"></i>${guide.views} views</span>
                    <span>${window.helpers.formatDate(guide.created_at)}</span>
                </div>
                <div class="flex gap-2">
                    <button onclick="app.viewItem('guide', ${JSON.stringify(guide).replace(/"/g, '&quot;')})" 
                            class="flex-1 px-3 py-2 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
                            <i class="fas fa-eye mr-1"></i> Xem
                    </button>
                    <button onclick="app.editItem('guide', ${JSON.stringify(guide).replace(/"/g, '&quot;')})" 
                            class="px-3 py-2 text-sm bg-green-500 text-white rounded hover:bg-green-600 transition-colors">
                            <i class="fas fa-edit"></i>
                    </button>
                    <button onclick="app.deleteItem('guides', ${guide.id})" 
                            class="px-3 py-2 text-sm bg-red-500 text-white rounded hover:bg-red-600 transition-colors">
                            <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        </div>
    `;
}

function formatNotificationRow(notification, app) {
    const readBadge = notification.is_read ? 'badge-blue' : 'badge-yellow';
    const sentBadge = notification.sent ? 'badge-green' : 'badge-red';
    
    return `
        <tr class="table-row-hover">
            <td class="px-6 py-4">
                <div class="flex flex-col gap-2">
                    <span class="px-2 py-1 text-xs font-medium rounded-full inline-block w-fit ${readBadge}">
                        ${notification.is_read ? '✓ Đã đọc' : '○ Chưa đọc'}
                    </span>
                    <span class="px-2 py-1 text-xs font-medium rounded-full inline-block w-fit ${sentBadge}">
                        ${notification.sent ? '✓ Đã gửi' : '⏳ Chưa gửi'}
                    </span>
                </div>
            </td>
            <td class="px-6 py-4 text-sm font-medium text-gray-900">${notification.title}</td>
            <td class="px-6 py-4 text-sm text-gray-600 line-clamp-2 max-w-xs">${notification.body}</td>
            <td class="px-6 py-4">
                <span class="px-2 py-1 text-xs font-medium rounded-full badge-purple">${notification.type}</span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-600">${window.helpers.getUserName(app.users, notification.user_id)}</td>
            <td class="px-6 py-4 text-sm space-x-3">
                <button onclick="app.viewItem('notification', ${JSON.stringify(notification).replace(/"/g, '&quot;')})" 
                        class="text-blue-600 hover:text-blue-800 transition-colors">
                        <i class="fas fa-eye"></i>
                </button>
                <button onclick="app.editItem('notification', ${JSON.stringify(notification).replace(/"/g, '&quot;')})" 
                        class="text-green-600 hover:text-green-800 transition-colors">
                        <i class="fas fa-edit"></i>
                </button>
                <button onclick="app.deleteItem('notifications', ${notification.id})" 
                        class="text-red-600 hover:text-red-800 transition-colors">
                        <i class="fas fa-trash"></i>
                </button>
                ${!notification.sent ? `
                    <button onclick="app.sendNotification(${notification.id})" 
                            class="text-purple-600 hover:text-purple-800 transition-colors"
                            title="Send">
                            <i class="fas fa-paper-plane"></i>
                    </button>
                ` : ''}
            </td>
        </tr>
    `;
}

// Export formatters
window.formatters = {
    formatUserRow,
    formatChecklistRow,
    formatLocationRow,
    formatGuideCard,
    formatNotificationRow
};