// Filter functions for data

function applyGlobalSearch(data, searchTerm, currentPage, users) {
    if (!searchTerm) return data;
    
    const term = searchTerm.toLowerCase();
    
    switch (currentPage) {
        case 'users':
            return data.filter(item => 
                item.name.toLowerCase().includes(term) ||
                item.email.toLowerCase().includes(term)
            );
            
        case 'checklist':
            return data.filter(item => 
                item.title.toLowerCase().includes(term) ||
                (item.description && item.description.toLowerCase().includes(term)) ||
                item.category.toLowerCase().includes(term)
            );
            
        case 'locations':
            return data.filter(item => {
                const userName = window.helpers.getUserName(users, item.user_id).toLowerCase();
                return userName.includes(term) ||
                       item.lat.toString().includes(term) ||
                       item.lon.toString().includes(term);
            });
            
        case 'guides':
            return data.filter(item => 
                item.title.toLowerCase().includes(term) ||
                item.content.toLowerCase().includes(term) ||
                item.category.toLowerCase().includes(term)
            );
            
        case 'notifications':
            return data.filter(item => 
                item.title.toLowerCase().includes(term) ||
                item.body.toLowerCase().includes(term) ||
                item.type.toLowerCase().includes(term)
            );
            
        default:
            return data;
    }
}

function applyChecklistFilters(data, filters, users) {
    let filtered = [...data];
    
    if (filters.user_id) {
        filtered = filtered.filter(item => item.user_id == filters.user_id);
    }
    
    if (filters.category) {
        filtered = filtered.filter(item => item.category === filters.category);
    }
    
    if (filters.status) {
        if (filters.status === 'checked') {
            filtered = filtered.filter(item => item.is_checked);
        } else if (filters.status === 'unchecked') {
            filtered = filtered.filter(item => !item.is_checked);
        }
    }
    
    return filtered;
}

function applyLocationsFilters(data, filters) {
    if (filters.user_id) {
        return data.filter(item => item.user_id == filters.user_id);
    }
    return data;
}

function applyGuidesFilters(data, filters) {
    let filtered = [...data];
    
    if (filters.category) {
        filtered = filtered.filter(item => item.category === filters.category);
    }
    
    if (filters.difficulty) {
        filtered = filtered.filter(item => item.difficulty === filters.difficulty);
    }
    
    if (filters.is_read !== '') {
        filtered = filtered.filter(item => item.is_read === (filters.is_read === 'true'));
    }
    
    return filtered;
}

function applyNotificationsFilters(data, filters) {
    let filtered = [...data];
    
    if (filters.user_id) {
        filtered = filtered.filter(item => item.user_id == filters.user_id);
    }
    
    if (filters.type) {
        filtered = filtered.filter(item => item.type === filters.type);
    }
    
    if (filters.is_read !== '') {
        filtered = filtered.filter(item => item.is_read === (filters.is_read === 'true'));
    }
    
    if (filters.sent !== '') {
        filtered = filtered.filter(item => item.sent === (filters.sent === 'true'));
    }
    
    return filtered;
}

// Export filter functions
window.filters = {
    applyGlobalSearch,
    applyChecklistFilters,
    applyLocationsFilters,
    applyGuidesFilters,
    applyNotificationsFilters
};