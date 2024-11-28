// Initialize components
const auth = new Auth();
const menu = new Menu();
const cart = new Cart();
const orders = new Orders();

// Handle navigation
document.querySelectorAll('.nav-link[data-page]').forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const page = e.target.dataset.page;
        
        // Update active link
        document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
        e.target.classList.add('active');

        // Load page content
        switch (page) {
            case 'menu':
                menu.loadMenu();
                break;
            case 'orders':
                if (!localStorage.getItem('token')) {
                    auth.showAuthModal('login');
                    return;
                }
                orders.loadOrders();
                break;
        }
    });
});