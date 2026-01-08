// Checklist Page Component

function ChecklistPage(app) {
    const items = app.filteredData;
    
    // Category colors
    const categoryColors = {
        'supplies': 'bg-blue-100 text-blue-800',
        'documents': 'bg-purple-100 text-purple-800',
        'emergency': 'bg-red-100 text-red-800',
        'food': 'bg-green-100 text-green-800',
        'water': 'bg-cyan-100 text-cyan-800',
        'shelter': 'bg-orange-100 text-orange-800'
    };
    
    return `
        <div class="space-y-6">
            <!-- Stats Cards -->
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Total Items</p>
                            <p class="text-2xl font-bold text-gray-800">${app.checklist.length}</p>
                        </div>
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-list text-blue-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Completed</p>
                            <p class="text-2xl font-bold text-green-600">${app.checklist.filter(c => c.is_checked).length}</p>
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
                            <p class="text-2xl font-bold text-yellow-600">${app.checklist.filter(c => !c.is_checked).length}</p>
                        </div>
                        <div class="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-clock text-yellow-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Progress</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${app.checklist.length > 0 ? Math.round((app.checklist.filter(c => c.is_checked).length / app.checklist.length) * 100) : 0}%
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-chart-pie text-purple-600 text-xl"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Filters -->
            <div class="bg-white rounded-lg shadow p-4">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">User</label>
                        <select x-model="filters.checklist.user_id" @change="applyFilters()" 
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="">All Users</option>
                            ${app.users.map(user => `
                                <option value="${user.id}">${user.name}</option>
                            `).join('')}
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
                        <select x-model="filters.checklist.category" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="">All Categories</option>
                            <option value="supplies">Supplies</option>
                            <option value="documents">Documents</option>
                            <option value="emergency">Emergency</option>
                            <option value="food">Food</option>
                            <option value="water">Water</option>
                            <option value="shelter">Shelter</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
                        <select x-model="filters.checklist.status" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                            <option value="">All Status</option>
                            <option value="checked">Completed</option>
                            <option value="unchecked">Pending</option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Checklist Items Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                ${items.length > 0 ? items.map(item => `
                    <div class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow duration-200">
                        <div class="p-6">
                            <div class="flex items-start justify-between mb-4">
                                <div class="flex items-center">
                                    <input type="checkbox" 
                                           ${item.is_checked ? 'checked' : ''} 
                                           @click="toggleChecked(${item.id})"
                                           class="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500 cursor-pointer">
                                    <h3 class="ml-3 text-lg font-semibold ${item.is_checked ? 'text-gray-400 line-through' : 'text-gray-800'}">
                                        ${item.title}
                                    </h3>
                                </div>
                            </div>
                            
                            <p class="text-sm text-gray-600 mb-4">
                                ${item.description || 'No description'}
                            </p>
                            
                            <div class="flex items-center justify-between">
                                <div class="flex items-center gap-2">
                                    <span class="px-2 py-1 text-xs font-semibold rounded-full ${categoryColors[item.category] || 'bg-gray-100 text-gray-800'}">
                                        ${item.category}
                                    </span>
                                    <span class="px-2 py-1 text-xs font-semibold rounded-full ${item.is_checked ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}">
                                        ${item.is_checked ? 'Done' : 'Pending'}
                                    </span>
                                </div>
                            </div>
                            
                            <div class="mt-4 pt-4 border-t border-gray-200 flex items-center justify-between">
                                <div class="flex items-center text-sm text-gray-500">
                                    <i class="fas fa-user mr-2"></i>
                                    ${window.helpers ? window.helpers.getUserName(app.users, item.user_id) : 'User ' + item.user_id}
                                </div>
                                <div class="flex items-center gap-2">
                                    <button @click="viewItem('checklist', ${JSON.stringify(item).replace(/"/g, '&quot;')})" 
                                            class="text-blue-600 hover:text-blue-900">
                                        <i class="fas fa-eye"></i>
                                    </button>
                                    <button @click="editItem('checklist', ${JSON.stringify(item).replace(/"/g, '&quot;')})" 
                                            class="text-indigo-600 hover:text-indigo-900">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button @click="deleteItem('checklist', ${item.id})" 
                                            class="text-red-600 hover:text-red-900">
                                        <i class="fas fa-trash"></i>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                `).join('') : `
                    <div class="col-span-full bg-white rounded-lg shadow p-12 text-center">
                        <i class="fas fa-clipboard-list text-6xl text-gray-300 mb-4"></i>
                        <p class="text-xl text-gray-500">No checklist items found</p>
                        <p class="text-sm text-gray-400 mt-2">Click "Thêm mới" to create your first item</p>
                    </div>
                `}
            </div>
        </div>
    `;
}

// Export to window
window.ChecklistPage = ChecklistPage;