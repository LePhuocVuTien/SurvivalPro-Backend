// Toast component

function ToastComponent(app) {
    return {
        getToastTemplate() {
            const icon = app.toast.type === 'success' ? 'fas fa-check-circle' : 'fas fa-exclamation-circle';
            return `
                <i class="${icon} text-white mr-3"></i>
                <p class="text-white font-medium">${app.toast.message}</p>
                <button onclick="app.toast.show = false" class="ml-4 text-white hover:text-gray-200">
                    <i class="fas fa-times"></i>
                </button>
            `;
        }
    };
}

// Export toast component
window.ToastComponent = ToastComponent;