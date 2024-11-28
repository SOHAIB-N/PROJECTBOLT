class Auth {
    constructor() {
        this.loginForm = document.querySelector('#loginForm form');
        this.registerForm = document.querySelector('#registerForm form');
        this.loginBtn = document.getElementById('loginBtn');
        this.registerBtn = document.getElementById('registerBtn');
        this.logoutBtn = document.getElementById('logoutBtn');
        this.authModal = document.getElementById('authModal');
        this.closeBtn = document.querySelector('.modal .close');

        this.setupEventListeners();
        this.checkAuthStatus();
    }

    setupEventListeners() {
        this.loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        this.registerForm.addEventListener('submit', (e) => this.handleRegister(e));
        this.loginBtn.addEventListener('click', () => this.showAuthModal('login'));
        this.registerBtn.addEventListener('click', () => this.showAuthModal('register'));
        this.logoutBtn.addEventListener('click', () => this.handleLogout());
        this.closeBtn.addEventListener('click', () => this.hideAuthModal());
    }

    showAuthModal(type) {
        this.authModal.style.display = 'block';
        document.getElementById('loginForm').classList.toggle('hidden', type !== 'login');
        document.getElementById('registerForm').classList.toggle('hidden', type !== 'register');
    }

    hideAuthModal() {
        this.authModal.style.display = 'none';
    }

    async handleLogin(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        
        try {
            const response = await api.login({
                email: formData.get('email'),
                password: formData.get('password')
            });

            this.setAuthToken(response.token);
            this.hideAuthModal();
            this.updateUIForAuthenticatedUser(response.user);
            window.location.reload();
        } catch (error) {
            alert(error.message);
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        
        try {
            const response = await api.register({
                name: formData.get('name'),
                email: formData.get('email'),
                password: formData.get('password'),
                phone: formData.get('phone')
            });

            this.setAuthToken(response.token);
            this.hideAuthModal();
            this.updateUIForAuthenticatedUser(response.user);
            window.location.reload();
        } catch (error) {
            alert(error.message);
        }
    }

    handleLogout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        this.updateUIForUnauthenticatedUser();
        window.location.reload();
    }

    setAuthToken(token) {
        localStorage.setItem('token', token);
    }

    checkAuthStatus() {
        const token = localStorage.getItem('token');
        const user = JSON.parse(localStorage.getItem('user'));
        
        if (token && user) {
            this.updateUIForAuthenticatedUser(user);
        } else {
            this.updateUIForUnauthenticatedUser();
        }
    }

    updateUIForAuthenticatedUser(user) {
        localStorage.setItem('user', JSON.stringify(user));
        this.loginBtn.classList.add('hidden');
        this.registerBtn.classList.add('hidden');
        this.logoutBtn.classList.remove('hidden');
        document.getElementById('cart-icon').classList.remove('hidden');
    }

    updateUIForUnauthenticatedUser() {
        localStorage.removeItem('user');
        this.loginBtn.classList.remove('hidden');
        this.registerBtn.classList.remove('hidden');
        this.logoutBtn.classList.add('hidden');
        document.getElementById('cart-icon').classList.add('hidden');
    }
}