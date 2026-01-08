function app() {
    return {
        // State
        users: [],
        checklist: [],
        locations: [],
        guides: [],
        notifications: [],
        
        currentPage: 'users',
        globalSearch: '',
        debounceTimer: null,
        
        // Sidebar template
        sidebarTemplate: '',
        
        // Computed properties for templates
        get modalContent() {
            try {
                if (!window.ModalComponent) return '<p>Loading modal...</p>';
                const modalComponent = window.ModalComponent(this);
                return modalComponent.getModalContent();
            } catch (error) {
                console.error('Modal content error:', error);
                return '<p>Error loading modal</p>';
            }
        },
        
        get toastTemplate() {
            try {
                if (!window.ToastComponent) return '';
                const toastComponent = window.ToastComponent(this);
                return toastComponent.getToastTemplate();
            } catch (error) {
                console.error('Toast template error:', error);
                return '<span>Notification</span>';
            }
        },
        
        get pageContent() {
            try {
                return this.getPageContent();
            } catch (error) {
                console.error('Page content error:', error);
                return '<div class="p-8 text-center text-red-500">Error loading page</div>';
            }
        },
        
        filters: {
            checklist: {
                user_id: '',
                category: '',
                status: ''
            },
            locations: {
                user_id: ''
            },
            guides: {
                category: '',
                difficulty: '',
                is_read: ''
            },
            notifications: {
                user_id: '',
                type: '',
                is_read: '',
                sent: ''
            }
        },
        
        modal: {
            show: false,
            mode: 'create',
            type: 'user',
            title: '',
            data: {}
        },
        
        toast: {
            show: false,
            type: 'success',
            message: '',
            timeout: null
        },
        
        map: null,
        markers: [],
        detailMap: null,
        
        navigation: [
            { id: 'users', label: 'Users', icon: 'fas fa-users' },
            { id: 'checklist', label: 'Checklist Items', icon: 'fas fa-list-check' },
            { id: 'locations', label: 'User Locations', icon: 'fas fa-map-marker-alt' },
            { id: 'guides', label: 'Survival Guides', icon: 'fas fa-book' },
            { id: 'notifications', label: 'Notifications', icon: 'fas fa-bell' }
        ],
        
        // Computed
        get filteredData() {
            let data = this.getCurrentData();
            
            if (this.globalSearch) {
                const searchTerm = this.globalSearch.toLowerCase();
                data = data.filter(item => {
                    if (this.currentPage === 'users') {
                        return item.name.toLowerCase().includes(searchTerm) ||
                               item.email.toLowerCase().includes(searchTerm);
                    } else if (this.currentPage === 'checklist') {
                        return item.title.toLowerCase().includes(searchTerm) ||
                               (item.description && item.description.toLowerCase().includes(searchTerm)) ||
                               item.category.toLowerCase().includes(searchTerm);
                    } else if (this.currentPage === 'locations') {
                        const userName = this.getUserName(item.user_id).toLowerCase();
                        return userName.includes(searchTerm) ||
                               item.lat.toString().includes(searchTerm) ||
                               item.lon.toString().includes(searchTerm);
                    } else if (this.currentPage === 'guides') {
                        return item.title.toLowerCase().includes(searchTerm) ||
                               item.content.toLowerCase().includes(searchTerm) ||
                               item.category.toLowerCase().includes(searchTerm);
                    } else if (this.currentPage === 'notifications') {
                        return item.title.toLowerCase().includes(searchTerm) ||
                               item.body.toLowerCase().includes(searchTerm) ||
                               item.type.toLowerCase().includes(searchTerm);
                    }
                    return true;
                });
            }
            
            // Apply specific filters
            if (this.currentPage === 'checklist') {
                if (this.filters.checklist.user_id) {
                    data = data.filter(item => item.user_id == this.filters.checklist.user_id);
                }
                if (this.filters.checklist.category) {
                    data = data.filter(item => item.category === this.filters.checklist.category);
                }
                if (this.filters.checklist.status) {
                    if (this.filters.checklist.status === 'checked') {
                        data = data.filter(item => item.is_checked);
                    } else if (this.filters.checklist.status === 'unchecked') {
                        data = data.filter(item => !item.is_checked);
                    }
                }
            } else if (this.currentPage === 'locations') {
                if (this.filters.locations.user_id) {
                    data = data.filter(item => item.user_id == this.filters.locations.user_id);
                }
            } else if (this.currentPage === 'guides') {
                if (this.filters.guides.category) {
                    data = data.filter(item => item.category === this.filters.guides.category);
                }
                if (this.filters.guides.difficulty) {
                    data = data.filter(item => item.difficulty === this.filters.guides.difficulty);
                }
                if (this.filters.guides.is_read !== '') {
                    data = data.filter(item => item.is_read === (this.filters.guides.is_read === 'true'));
                }
            } else if (this.currentPage === 'notifications') {
                if (this.filters.notifications.user_id) {
                    data = data.filter(item => item.user_id == this.filters.notifications.user_id);
                }
                if (this.filters.notifications.type) {
                    data = data.filter(item => item.type === this.filters.notifications.type);
                }
                if (this.filters.notifications.is_read !== '') {
                    data = data.filter(item => item.is_read === (this.filters.notifications.is_read === 'true'));
                }
                if (this.filters.notifications.sent !== '') {
                    data = data.filter(item => item.sent === (this.filters.notifications.sent === 'true'));
                }
            }
            
            return data;
        },
        
        // Lifecycle
        init() {
            console.log('🚀 App init started');
            
            try {
                // Generate mock data first
                this.generateMockData();
                console.log('✅ Mock data generated:', {
                    users: this.users.length,
                    checklist: this.checklist.length,
                    locations: this.locations.length,
                    guides: this.guides.length,
                    notifications: this.notifications.length
                });
                
                // Update sidebar
                this.updateSidebar();
                console.log('✅ Sidebar updated');
                
                // Render page
                this.renderPage();
                console.log('✅ Page rendered');
                
                console.log('🎉 App initialized successfully');
            } catch (error) {
                console.error('❌ App init error:', error);
            }
        },
        
        // Methods
        generateMockData() {
            try {
                if (typeof window.getMockData !== 'function') {
                    console.error('❌ getMockData is not defined');
                    // Fallback to empty data
                    this.users = [];
                    this.checklist = [];
                    this.locations = [];
                    this.guides = [];
                    this.notifications = [];
                    return;
                }
                
                const data = window.getMockData();
                this.users = data.users || [];
                this.checklist = data.checklist || [];
                this.locations = data.locations || [];
                this.guides = data.guides || [];
                this.notifications = data.notifications || [];
                
                console.log('✅ Data loaded:', data);
            } catch (error) {
                console.error('❌ Error generating mock data:', error);
            }
        },
        
        // Sidebar methods
        updateSidebar() {
            try {
                if (!window.SidebarComponent) {
                    console.warn('⚠️ SidebarComponent not loaded yet');
                    this.sidebarTemplate = '<div class="p-4 text-gray-500">Loading sidebar...</div>';
                    return;
                }
                const sidebarComponent = window.SidebarComponent(this);
                this.sidebarTemplate = sidebarComponent.getSidebarTemplate();
                console.log('✅ Sidebar template updated');
            } catch (error) {
                console.error('❌ Error updating sidebar:', error);
                this.sidebarTemplate = '<div class="p-4 text-red-500">Error loading sidebar</div>';
            }
        },
        
        getCurrentPageTitle() {
            const page = this.navigation.find(item => item.id === this.currentPage);
            return page ? page.label : 'Dashboard';
        },
        
        getCurrentData() {
            switch (this.currentPage) {
                case 'users': return this.users;
                case 'checklist': return this.checklist;
                case 'locations': return this.locations;
                case 'guides': return this.guides;
                case 'notifications': return this.notifications;
                default: return [];
            }
        },
        
        changePage(page) {
            this.currentPage = page;
            this.globalSearch = '';
            this.resetFilters();
            // Only update sidebar when page changes
            this.$nextTick(() => {
                this.updateSidebar();
            });
        },
        
        resetFilters() {
            this.filters = {
                checklist: {
                    user_id: '',
                    category: '',
                    status: ''
                },
                locations: {
                    user_id: ''
                },
                guides: {
                    category: '',
                    difficulty: '',
                    is_read: ''
                },
                notifications: {
                    user_id: '',
                    type: '',
                    is_read: '',
                    sent: ''
                }
            };
        },
        
        applyFilters() {
            // Được gọi khi filter thay đổi, nhưng computed property filteredData đã xử lý
        },
        
        debounceSearch() {
            clearTimeout(this.debounceTimer);
            this.debounceTimer = setTimeout(() => {
                this.applyFilters();
            }, 300);
        },
        
        getUserName(userId) {
            const user = this.users.find(u => u.id == userId);
            return user ? user.name : `User ${userId}`;
        },
        
        formatDate(dateString) {
            if (!dateString) return 'N/A';
            const date = new Date(dateString);
            return date.toLocaleDateString('vi-VN') + ' ' + date.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
        },
        
        openCreateModal() {
            let emptyData = {};
            if (this.currentPage === 'users') {
                emptyData = { email: '', password: '', name: '', avatar: '', push_token: '' };
                this.modal.type = 'user';
                this.modal.title = 'Thêm User mới';
            } else if (this.currentPage === 'checklist') {
                emptyData = { user_id: '', title: '', category: 'supplies', description: '', is_checked: false };
                this.modal.type = 'checklist';
                this.modal.title = 'Thêm Checklist Item mới';
            } else if (this.currentPage === 'guides') {
                emptyData = { user_id: '', title: '', category: 'first-aid', difficulty: 'easy', icon: '', content: '', image_url: '', views: 0, is_read: false };
                this.modal.type = 'guide';
                this.modal.title = 'Thêm Survival Guide mới';
            } else if (this.currentPage === 'notifications') {
                emptyData = { user_id: '', title: '', body: '', data: '', type: 'push', is_read: false, sent: false };
                this.modal.type = 'notification';
                this.modal.title = 'Thêm Notification mới';
            }
            
            this.modal.mode = 'create';
            this.modal.data = { ...emptyData };
            this.modal.show = true;
        },
        
        viewItem(type, item) {
            this.modal.type = type;
            this.modal.mode = 'view';
            this.modal.data = { ...item };
            this.modal.title = `Chi tiết ${this.getModalTypeName(type)}`;
            this.modal.show = true;
        },
        
        editItem(type, item) {
            this.modal.type = type;
            this.modal.mode = 'edit';
            this.modal.data = { ...item };
            this.modal.title = `Sửa ${this.getModalTypeName(type)}`;
            this.modal.show = true;
        },
        
        getModalTypeName(type) {
            const names = {
                'user': 'User',
                'checklist': 'Checklist Item',
                'location': 'Location',
                'guide': 'Survival Guide',
                'notification': 'Notification'
            };
            return names[type] || 'Item';
        },
        
        closeModal() {
            this.modal.show = false;
            this.modal.data = {};
        },
        
        saveItem() {
            // Generate new ID for create mode
            if (this.modal.mode === 'create') {
                const dataArray = this.getCurrentData();
                const maxId = dataArray.length > 0 ? Math.max(...dataArray.map(item => item.id)) : 0;
                this.modal.data.id = maxId + 1;
                this.modal.data.created_at = new Date().toISOString();
                
                // Add to appropriate array
                if (this.currentPage === 'users') {
                    this.users.push({ ...this.modal.data });
                } else if (this.currentPage === 'checklist') {
                    this.checklist.push({ ...this.modal.data });
                } else if (this.currentPage === 'guides') {
                    this.guides.push({ ...this.modal.data });
                } else if (this.currentPage === 'notifications') {
                    this.notifications.push({ ...this.modal.data });
                }
                
                this.showToast('Tạo mới thành công!', 'success');
            } else if (this.modal.mode === 'edit') {
                // Update existing item
                const dataArray = this.getCurrentData();
                const index = dataArray.findIndex(item => item.id === this.modal.data.id);
                if (index !== -1) {
                    dataArray[index] = { ...this.modal.data };
                    this.showToast('Cập nhật thành công!', 'success');
                }
            }
            
            this.closeModal();
            // Update sidebar after data change
            this.$nextTick(() => {
                this.updateSidebar();
            });
        },
        
        deleteItem(type, id) {
            if (!confirm('Bạn có chắc chắn muốn xóa mục này?')) return;
            
            let dataArray;
            switch (type) {
                case 'users': dataArray = this.users; break;
                case 'checklist': dataArray = this.checklist; break;
                case 'guides': dataArray = this.guides; break;
                case 'notifications': dataArray = this.notifications; break;
                default: return;
            }
            
            const index = dataArray.findIndex(item => item.id === id);
            if (index !== -1) {
                dataArray.splice(index, 1);
                this.showToast('Xóa thành công!', 'success');
                // Update sidebar after delete
                this.$nextTick(() => {
                    this.updateSidebar();
                });
            }
        },
        
        toggleChecked(id) {
            const item = this.checklist.find(item => item.id === id);
            if (item) {
                item.is_checked = !item.is_checked;
                this.showToast(`Đã ${item.is_checked ? 'đánh dấu hoàn thành' : 'bỏ đánh dấu'}!`, 'success');
                // Update sidebar after toggle
                this.$nextTick(() => {
                    this.updateSidebar();
                });
            }
        },
        
        sendNotification(id) {
            const notification = this.notifications.find(n => n.id === id);
            if (notification) {
                notification.sent = true;
                notification.sent_at = new Date().toISOString();
                this.showToast('Đã gửi thông báo thành công!', 'success');
                // Update sidebar after send
                this.$nextTick(() => {
                    this.updateSidebar();
                });
            }
        },
        
        showToast(message, type = 'success') {
            this.toast.message = message;
            this.toast.type = type;
            this.toast.show = true;
            
            clearTimeout(this.toast.timeout);
            this.toast.timeout = setTimeout(() => {
                this.toast.show = false;
            }, 3000);
        },
        
        renderPage() {
            // Don't call updateSidebar here to aclvoid infinite loop
            console.log('📄 Rendering page:', this.currentPage);
        },
        
        getPageContent() {
            try {
                switch (this.currentPage) {
                    case 'users':
                        return window.UsersPage ? window.UsersPage(this) : '<div class="p-8">Loading Users...</div>';
                    case 'checklist':
                        return window.ChecklistPage ? window.ChecklistPage(this) : '<div class="p-8">Loading Checklist...</div>';
                    case 'locations':
                        setTimeout(() => {
                            if (typeof window.MapComponent !== 'undefined') {
                                window.MapComponent.initMap(this);
                            }
                        }, 100);
                        return window.LocationsPage ? window.LocationsPage(this) : '<div class="p-8">Loading Locations...</div>';
                    case 'guides':
                        return window.GuidesPage ? window.GuidesPage(this) : '<div class="p-8">Loading Guides...</div>';
                    case 'notifications':
                        return window.NotificationsPage ? window.NotificationsPage(this) : '<div class="p-8">Loading Notifications...</div>';
                    default:
                        return '<div class="text-center text-gray-500 p-8">Trang không tồn tại</div>';
                }
            } catch (error) {
                console.error('❌ Error getting page content:', error);
                return '<div class="text-center text-red-500 p-8">Error loading page</div>';
            }
        }
    };
}

// Khởi tạo ứng dụng
document.addEventListener('alpine:init', () => {
    console.log('🚀 Alpine init started');
    Alpine.data('app', app);
    console.log('✅ App registered with Alpine');
});

// Debug: Log khi Alpine ready
document.addEventListener('alpine:initialized', () => {
    console.log('✅ Alpine fully initialized');
});