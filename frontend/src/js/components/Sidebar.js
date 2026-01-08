// Sidebar component

function SidebarComponent(app) {
    // Sidebar menu items configuration
    const menuItems = [
        {
            id: 'users',
            title: 'Users',
            icon: 'fa-users',
            badge: () => (app.users || []).length,
            badgeColor: 'bg-blue-500'
        },
        {
            id: 'checklist',
            title: 'Checklists',
            icon: 'fa-clipboard-list',
            badge: () => (app.checklist || []).length,
            badgeColor: 'bg-green-500'
        },
        {
            id: 'guides',
            title: 'Survival Guides',
            icon: 'fa-book',
            badge: () => (app.guides || []).length,
            badgeColor: 'bg-purple-500'
        },
        {
            id: 'locations',
            title: 'Locations',
            icon: 'fa-map-marker-alt',
            badge: () => (app.locations || []).length,
            badgeColor: 'bg-red-500'
        },
        {
            id: 'notifications',
            title: 'Notifications',
            icon: 'fa-bell',
            badge: () => {
                const unread = (app.notifications || []).filter(n => !n.is_read).length;
                return unread > 0 ? unread : null;
            },
            badgeColor: 'bg-orange-500'
        }
    ];

    return {
        /**
         * Generate sidebar HTML template
         */
        getSidebarTemplate() {
            return `
                <div class="space-y-1">
                    ${menuItems.map(item => this.generateMenuItem(item)).join('')}
                </div>
                
                <div class="mt-8 pt-4 border-t border-gray-200">
                    <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3 px-3">
                        Statistics
                    </div>
                    ${this.generateStatsSection()}
                </div>
            `;
        },

        /**
         * Generate a single menu item
         */
        generateMenuItem(item) {
            const isActive = app.currentPage === item.id;
            const badge = typeof item.badge === 'function' ? item.badge() : item.badge;
            
            return `
                <button @click="changePage('${item.id}')" 
                        class="sidebar-item ${isActive ? 'active' : ''} w-full flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200 group">
                    <div class="flex items-center gap-3">
                        <i class="fas ${item.icon} text-lg ${isActive ? 'text-blue-600' : 'text-gray-500 group-hover:text-blue-600'}"></i>
                        <span class="font-medium ${isActive ? 'text-blue-600' : 'text-gray-700 group-hover:text-blue-600'}">
                            ${item.title}
                        </span>
                    </div>
                    ${badge ? `
                        <span class="${item.badgeColor} text-white text-xs font-bold px-2 py-1 rounded-full min-w-[24px] text-center">
                            ${badge}
                        </span>
                    ` : ''}
                </button>
            `;
        },

        /**
         * Generate statistics section
         */
        generateStatsSection() {
            const stats = [
                {
                    label: 'Total Users',
                    value: (app.users || []).length,
                    icon: 'fa-user',
                    color: 'text-blue-600'
                },
                {
                    label: 'Active Checklists',
                    value: (app.checklist || []).filter(c => !c.is_checked).length,
                    icon: 'fa-tasks',
                    color: 'text-green-600'
                },
                {
                    label: 'Unread Guides',
                    value: (app.guides || []).filter(g => !g.is_read).length,
                    icon: 'fa-book-open',
                    color: 'text-purple-600'
                },
                {
                    label: 'Pending Notifications',
                    value: (app.notifications || []).filter(n => !n.sent).length,
                    icon: 'fa-paper-plane',
                    color: 'text-orange-600'
                }
            ];

            return `
                <div class="space-y-3">
                    ${stats.map(stat => `
                        <div class="px-3 py-2 bg-gray-50 rounded-lg">
                            <div class="flex items-center justify-between">
                                <div class="flex items-center gap-2">
                                    <i class="fas ${stat.icon} ${stat.color} text-sm"></i>
                                    <span class="text-xs text-gray-600">${stat.label}</span>
                                </div>
                                <span class="text-sm font-bold text-gray-800">${stat.value}</span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            `;
        },

        /**
         * Get current active menu item
         */
        getActiveMenuItem() {
            return menuItems.find(item => item.id === app.currentPage);
        },

        /**
         * Get menu item by ID
         */
        getMenuItem(id) {
            return menuItems.find(item => item.id === id);
        },

        /**
         * Get all menu items
         */
        getMenuItems() {
            return menuItems;
        }
    };
}

// Export sidebar component
window.SidebarComponent = SidebarComponent;