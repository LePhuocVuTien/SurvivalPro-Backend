// Locations Page Component

function LocationsPage(app) {
    const locations = app.filteredData;
    
    // Group locations by user
    const locationsByUser = {};
    locations.forEach(loc => {
        if (!locationsByUser[loc.user_id]) {
            locationsByUser[loc.user_id] = [];
        }
        locationsByUser[loc.user_id].push(loc);
    });
    
    return `
        <div class="space-y-6">
            <!-- Stats Cards -->
            <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Total Locations</p>
                            <p class="text-2xl font-bold text-gray-800">${app.locations.length}</p>
                        </div>
                        <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-map-marker-alt text-red-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Tracked Users</p>
                            <p class="text-2xl font-bold text-gray-800">${Object.keys(locationsByUser).length}</p>
                        </div>
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-users text-blue-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Latest Update</p>
                            <p class="text-sm font-bold text-gray-800">
                                ${locations.length > 0 ? new Date(Math.max(...locations.map(l => new Date(l.timestamp)))).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) : 'N/A'}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-clock text-green-600 text-xl"></i>
                        </div>
                    </div>
                </div>
                
                <div class="bg-white rounded-lg shadow p-6">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-sm text-gray-500">Active Today</p>
                            <p class="text-2xl font-bold text-gray-800">
                                ${locations.filter(l => {
                                    const today = new Date();
                                    const locDate = new Date(l.timestamp);
                                    return locDate.toDateString() === today.toDateString();
                                }).length}
                            </p>
                        </div>
                        <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
                            <i class="fas fa-signal text-purple-600 text-xl"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Filter -->
            <div class="bg-white rounded-lg shadow p-4">
                <div class="flex items-center gap-4">
                    <div class="flex-1">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Filter by User</label>
                        <select x-model="filters.locations.user_id" @change="applyFilters()"
                                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500">
                            <option value="">All Users</option>
                            ${app.users.map(user => `
                                <option value="${user.id}">${user.name}</option>
                            `).join('')}
                        </select>
                    </div>
                    <div class="pt-7">
                        <button @click="if(window.MapComponent) window.MapComponent.initMap($data)" 
                                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors flex items-center gap-2">
                            <i class="fas fa-sync-alt"></i>
                            <span>Refresh Map</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Map -->
            <div class="bg-white rounded-lg shadow overflow-hidden">
                <div id="map" style="height: 500px; width: 100%;"></div>
            </div>

            <!-- Locations Timeline -->
            <div class="bg-white rounded-lg shadow">
                <div class="p-6 border-b border-gray-200">
                    <h3 class="text-lg font-bold text-gray-800 flex items-center">
                        <i class="fas fa-route mr-2 text-red-600"></i>
                        Location Timeline
                    </h3>
                </div>
                
                <div class="p-6">
                    ${Object.keys(locationsByUser).length > 0 ? Object.keys(locationsByUser).map(userId => {
                        const user = app.users.find(u => u.id == userId);
                        const userLocs = locationsByUser[userId].sort((a, b) => 
                            new Date(b.timestamp) - new Date(a.timestamp)
                        );
                        
                        return `
                            <div class="mb-8 last:mb-0">
                                <div class="flex items-center gap-3 mb-4">
                                    <img src="${user?.avatar || 'https://via.placeholder.com/150'}" 
                                         class="w-12 h-12 rounded-full object-cover"
                                         onerror="this.src='https://via.placeholder.com/150'">
                                    <div>
                                        <h4 class="font-bold text-gray-800">${user?.name || 'User ' + userId}</h4>
                                        <p class="text-sm text-gray-500">${userLocs.length} location${userLocs.length > 1 ? 's' : ''}</p>
                                    </div>
                                </div>
                                
                                <div class="ml-6 pl-6 border-l-2 border-gray-200 space-y-4">
                                    ${userLocs.map((loc, index) => `
                                        <div class="relative">
                                            <div class="absolute -left-8 mt-1 w-4 h-4 rounded-full bg-red-600 border-2 border-white"></div>
                                            <div class="bg-gray-50 rounded-lg p-4 hover:bg-gray-100 transition-colors">
                                                <div class="flex items-start justify-between">
                                                    <div class="flex-1">
                                                        <div class="flex items-center gap-2 mb-2">
                                                            <i class="fas fa-map-marker-alt text-red-600"></i>
                                                            <span class="font-medium text-gray-800">
                                                                Location #${loc.id}
                                                            </span>
                                                            ${index === 0 ? '<span class="px-2 py-1 text-xs font-semibold bg-green-100 text-green-800 rounded-full">Latest</span>' : ''}
                                                        </div>
                                                        <div class="grid grid-cols-2 gap-4 text-sm">
                                                            <div>
                                                                <span class="text-gray-500">Latitude:</span>
                                                                <span class="font-mono text-gray-800 ml-2">${loc.latitude.toFixed(6)}</span>
                                                            </div>
                                                            <div>
                                                                <span class="text-gray-500">Longitude:</span>
                                                                <span class="font-mono text-gray-800 ml-2">${loc.longtitude.toFixed(6)}</span>
                                                            </div>
                                                        </div>
                                                        <div class="mt-2 text-xs text-gray-500">
                                                            <i class="fas fa-clock mr-1"></i>
                                                            ${window.helpers.formatDate(loc.timestamp)}
                                                        </div>
                                                    </div>
                                                    <div class="flex items-center gap-2">
                                                        <button @click="alert('View on map: ${loc.latitude}, ${loc.longtitude}')" 
                                                                class="p-2 text-blue-600 hover:text-blue-900">
                                                            <i class="fas fa-eye"></i>
                                                        </button>
                                                        <a href="https://www.google.com/maps?q=${loc.latitude},${loc.longtitude}" 
                                                           target="_blank"
                                                           class="p-2 text-green-600 hover:text-green-900">
                                                            <i class="fas fa-external-link-alt"></i>
                                                        </a>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    `).join('')}
                                </div>
                            </div>
                        `;
                    }).join('') : `
                        <div class="text-center py-12">
                            <i class="fas fa-map-marked-alt text-6xl text-gray-300 mb-4"></i>
                            <p class="text-xl text-gray-500">No location data available</p>
                        </div>
                    `}
                </div>
            </div>
        </div>
    `;
}

// Export to window
window.LocationsPage = LocationsPage;