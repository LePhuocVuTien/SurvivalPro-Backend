// Helper functions for the application

// Format date to Vietnamese format
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('vi-VN') + ' ' + date.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
    } catch (error) {
        return 'Invalid Date';
    }
}

// Debounce function
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Generate unique ID
function generateId(existingData) {
    if (!existingData || existingData.length === 0) return 1;
    const maxId = Math.max(...existingData.map(item => item.id));
    return maxId + 1;
}

// Get user name by ID
function getUserName(users, userId) {
    const user = users.find(u => u.id == userId);
    return user ? user.name : `User ${userId}`;
}

// Get modal type name
function getModalTypeName(type) {
    const names = {
        'user': 'User',
        'checklist': 'Checklist Item',
        'location': 'Location',
        'guide': 'Survival Guide',
        'notification': 'Notification'
    };
    return names[type] || 'Item';
}

// Get badge class by category
function getBadgeClassByCategory(category) {
    const classes = {
        'supplies': 'badge-blue',
        'documents': 'badge-purple',
        'emergency': 'badge-red',
        'food': 'badge-green',
        'water': 'badge-cyan',
        'shelter': 'badge-amber',
        'first-aid': 'badge-red',
        'navigation': 'badge-blue',
        'fire': 'badge-yellow',
        'signaling': 'badge-purple'
    };
    return classes[category] || 'badge-blue';
}

// Get difficulty badge class
function getDifficultyBadgeClass(difficulty) {
    const classes = {
        'easy': 'badge-green',
        'medium': 'badge-yellow',
        'hard': 'badge-red'
    };
    return classes[difficulty] || 'badge-blue';
}

// Validate email
function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

// Validate required fields
function validateRequired(fields, data) {
    for (const field of fields) {
        if (!data[field] || data[field].toString().trim() === '') {
            return { isValid: false, message: `${field} là bắt buộc` };
        }
    }
    return { isValid: true };
}

// Deep clone object
function deepClone(obj) {
    return JSON.parse(JSON.stringify(obj));
}

// Export helpers
window.helpers = {
    formatDate,
    debounce,
    generateId,
    getUserName,
    getModalTypeName,
    getBadgeClassByCategory,
    getDifficultyBadgeClass,
    validateEmail,
    validateRequired,
    deepClone
};