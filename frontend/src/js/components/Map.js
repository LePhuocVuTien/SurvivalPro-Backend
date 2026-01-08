// Map Component for Locations Page

window.MapComponent = {
    map: null,
    markers: [],
    
    /**
     * Initialize map
     */
    initMap(app) {
        try {
            // Wait for map container to be available
            setTimeout(() => {
                const mapContainer = document.getElementById('map');
                if (!mapContainer) {
                    console.warn('Map container not found');
                    return;
                }
                
                // Check if Leaflet is loaded
                if (typeof L === 'undefined') {
                    console.error('Leaflet library not loaded');
                    return;
                }
                
                // Initialize map if not already initialized
                if (!this.map) {
                    // Default center: Da Nang, Vietnam
                    this.map = L.map('map').setView([16.0544, 108.2022], 13);
                    
                    // Add tile layer
                    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                        attribution: '© OpenStreetMap contributors',
                        maxZoom: 19
                    }).addTo(this.map);
                    
                    console.log('✅ Map initialized');
                }
                
                // Clear existing markers
                this.clearMarkers();
                
                // Add markers for all locations
                if (app.locations && app.locations.length > 0) {
                    app.locations.forEach(location => {
                        this.addMarker(location, app);
                    });
                    
                    // Fit bounds to show all markers
                    if (this.markers.length > 0) {
                        const group = L.featureGroup(this.markers);
                        this.map.fitBounds(group.getBounds().pad(0.1));
                    }
                }
            }, 200);
        } catch (error) {
            console.error('Error initializing map:', error);
        }
    },
    
    /**
     * Add marker to map
     */
    addMarker(location, app) {
        try {
            const user = app.users.find(u => u.id === location.user_id);
            const userName = user ? user.name : `User ${location.user_id}`;
            
            const marker = L.marker([location.latitude, location.longtitude])
                .bindPopup(`
                    <div class="p-2">
                        <strong>${userName}</strong><br>
                        <small>${location.address || 'Unknown address'}</small><br>
                        <small class="text-gray-500">${window.helpers.formatDate(location.created_at)}</small>
                    </div>
                `)
                .addTo(this.map);
            
            this.markers.push(marker);
        } catch (error) {
            console.error('Error adding marker:', error);
        }
    },
    
    /**
     * Clear all markers
     */
    clearMarkers() {
        this.markers.forEach(marker => {
            if (this.map) {
                this.map.removeLayer(marker);
            }
        });
        this.markers = [];
    },
    
    /**
     * Destroy map instance
     */
    destroyMap() {
        if (this.map) {
            this.map.remove();
            this.map = null;
        }
        this.markers = [];
    }
};