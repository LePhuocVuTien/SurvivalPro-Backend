// Guides Page Component

function GuidesPage(app) {
    const guides = app.filteredData;
    
    // Category info
    const categoryInfo = {
        'first-aid': { name: 'First Aid', color: 'bg-red-100 text-red-800', icon: 'fa-medkit' },
        'shelter': { name: 'Shelter', color: 'bg-orange-100 text-orange-800', icon: 'fa-home' },
        'food': { name: 'Food & Water', color: 'bg-green-100 text-green-800', icon: 'fa-utensils' },
        'navigation': { name: 'Navigation', color: 'bg-blue-100 text-blue-800', icon: 'fa-compass' },
        'fire': { name: 'Fire Making', color: 'bg-yellow-100 text-yellow-800', icon: 'fa-fire' },
        'signaling': { name: 'Signaling', color: 'bg-purple-100 text-purple-800', icon: 'fa-broadcast-tower' }
    };
    
    // Difficulty colors
    const difficultyColors = {
        'easy': 'bg-green-100 text-green-800',
        'medium': 'bg-yellow-100 text-yellow-800',
        'hard': 'bg-red-100 text-red-800'
    };
    
    return `
        <div class="space-y-6">
            <!-- Stats Cards -->
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Total Guides</p>
                            <p class="text-2xl font-bold text-gray-800">${app.guides.length}</p>
                        </div>
                        <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-book text-purple-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Total Views</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${app.guides.reduce((sum, g) => sum + g.views, 0).toLocaleString()}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-eye text-blue-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Categories</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${new Set(app.guides.map(g => g.category)).size}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-layer-group text-green-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Most Popular</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${Math.max(...app.guides.map(g => g.views))}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-star text-orange-600 text-xl"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Filters -->
            <div class="bg-white rounded-lg shadow p-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
                        <select x-model="filters.guides.category" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500">
                            <option value="">All Categories</option>
                            <option value="first-aid">First Aid</option>
                            <option value="shelter">Shelter</option>
                            <option value="food">Food & Water</option>
                            <option value="navigation">Navigation</option>
                            <option value="fire">Fire Making</option>
                            <option value="signaling">Signaling</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Difficulty</label>
                        <select x-model="filters.guides.difficulty" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500">
                            <option value="">All Levels</option>
                            <option value="easy">Easy</option>
                            <option value="medium">Medium</option>
                            <option value="hard">Hard</option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Guides Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                ${guides.length > 0 ? guides.map(guide => {
                    const category = categoryInfo[guide.category] || { name: guide.category, color: 'bg-gray-100 text-gray-800', icon: 'fa-book' };
                    return `
                    <div class="bg-white rounded-lg shadow hover:shadow-xl transition-all duration-300 overflow-hidden">
                        <!-- Image -->
                        <div class="relative h-48 bg-gradient-to-br from-purple-400 to-indigo-600">
                            ${guide.image_url ? `
                                <img src="${guide.image_url}" 
                                     alt="${guide.title}" 
                                     class="w-full h-full object-cover"
                                     onerror="this.style.display='none'">
                            ` : ''}
                            <div class="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center">
                                <span class="text-6xl">${guide.icon}</span>
                            </div>
                            <div class="absolute top-4 right-4">
                                <span class="px-3 py-1 text-xs font-bold rounded-full ${difficultyColors[guide.difficulty] || 'bg-gray-100 text-gray-800'}">
                                    ${guide.difficulty.toUpperCase()}
                                </span>
                            </div>
                        </div>
                        
                        <!-- Content -->
                        <div class="p-6">
                            <div class="flex items-center gap-2 mb-3">
                                <span class="px-2 py-1 text-xs font-semibold rounded ${category.color}">
                                    <i class="fas ${category.icon} mr-1"></i>
                                    ${category.name}
                                </span>
                                <span class="text-xs text-gray-500 flex items-center">
                                    <i class="fas fa-eye mr-1"></i>
                                    ${guide.views.toLocaleString()}
                                </span>
                            </div>
                            
                            <h3 class="text-lg font-bold text-gray-800 mb-3 line-clamp-2">
                                ${guide.title}
                            </h3>
                            
                            <p class="text-sm text-gray-600 mb-4 line-clamp-3">
                                ${guide.content.split('\\n')[0]}
                            </p>
                            
                            <div class="flex items-center justify-between pt-4 border-t border-gray-200">
                                <div class="text-xs text-gray-500">
                                    <i class="fas fa-clock mr-1"></i>
                                    ${window.helpers ? new Date(guide.created_at).toLocaleDateString('vi-VN') : guide.created_at}
                                </div>
                                <div class="flex items-center gap-2">
                                    <button @click="viewItem('guide', ${JSON.stringify(guide).replace(/"/g, '&quot;').replace(/\n/g, '\\n')})" 
                                            class="px-3 py-1 text-xs font-medium text-white bg-purple-600 rounded hover:bg-purple-700 transition-colors">
                                        <i class="fas fa-book-open mr-1"></i>
                                        Read
                                    </button>
                                    <button @click="editItem('guide', ${JSON.stringify(guide).replace(/"/g, '&quot;').replace(/\n/g, '\\n')})" 
                                            class="p-2 text-indigo-600 hover:text-indigo-900">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button @click="deleteItem('guides', ${guide.id})" 
                                            class="p-2 text-red-600 hover:text-red-900">
                                        <i class="fas fa-trash"></i>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                    `;
                }).join('') : `
                    <div class="col-span-full bg-white rounded-lg shadow p-12 text-center">
                        <i class="fas fa-book text-6xl text-gray-300 mb-4"></i>
                        <p class="text-xl text-gray-500">No survival guides found</p>
                        <p class="text-sm text-gray-400 mt-2">Click "Thêm mới" to create a new guide</p>
                    </div>
                `}
            </div>
        </div>
        
        <style>
            .line-clamp-2 {
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }
            .line-clamp-3 {
                display: -webkit-box;
                -webkit-line-clamp: 3;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }
        </style>
    `;
}

// Export to window
window.GuidesPage = GuidesPage;