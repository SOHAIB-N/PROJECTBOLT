class Cart {
    constructor() {
        this.items = [];
        this.cartIcon = document.getElementById('cart-icon');
        this.cartModal = document.getElementById('cartModal');
        this.cartItems = document.getElementById('cartItems');
        this.cartTotal = document.getElementById('cartTotal');
        this.checkoutBtn = document.getElementById('checkoutBtn');
        this.closeBtn = this.cartModal.querySelector('.close');
        this.hostelDelivery = document.getElementById('hostelDelivery');
        this.roomNumber = document.getElementById('roomNumber');

        this.setupEventListeners();
    }

    setupEventListeners() {
        this.cartIcon.addEventListener('click', () => this.showCart());
        this.closeBtn.addEventListener('click', () => this.hideCart());
        this.checkoutBtn.addEventListener('click', () => this.checkout());
        this.hostelDelivery.addEventListener('change', () => {
            this.roomNumber.classList.toggle('hidden', !this.hostelDelivery.checked);
        });
    }

    addItem(item) {
        const existingItem = this.items.find(i => i.ID === item.ID);
        if (existingItem) {
            existingItem.quantity += 1;
        } else {
            this.items.push({ ...item, quantity: 1 });
        }
        this.updateCartUI();
    }

    removeItem(itemId) {
        this.items = this.items.filter(item => item.ID !== itemId);
        this.updateCartUI();
    }

    updateQuantity(itemId, quantity) {
        const item = this.items.find(i => i.ID === itemId);
        if (item) {
            item.quantity = Math.max(0, quantity);
            if (item.quantity === 0) {
                this.removeItem(itemId);
            }
        }
        this.updateCartUI();
    }

    updateCartUI() {
        const count = this.items.reduce((sum, item) => sum + item.quantity, 0);
        document.getElementById('cart-count').textContent = count;

        this.cartItems.innerHTML = this.items.map(item => `
            <div class="cart-item">
                <img src="${item.ImageURL || 'images/default-dish.jpg'}" alt="${item.Name}">
                <div class="cart-item-details">
                    <h4>${item.Name}</h4>
                    <p>$${item.Price.toFixed(2)}</p>
                    <div class="quantity-controls">
                        <button onclick="cart.updateQuantity(${item.ID}, ${item.quantity - 1})">-</button>
                        <span>${item.quantity}</span>
                        <button onclick="cart.updateQuantity(${item.ID}, ${item.quantity + 1})">+</button>
                    </div>
                </div>
                <button class="remove-item" onclick="cart.removeItem(${item.ID})">×</button>
            </div>
        `).join('');

        const total = this.items.reduce((sum, item) => sum + (item.Price * item.quantity), 0);
        this.cartTotal.innerHTML = `<h3>Total: $${total.toFixed(2)}</h3>`;
    }

    showCart() {
        this.cartModal.style.display = 'block';
        this.updateCartUI();
    }

    hideCart() {
        this.cartModal.style.display = 'none';
    }

    async checkout() {
        if (this.items.length === 0) {
            alert('Your cart is empty');
            return;
        }

        const isDelivery = this.hostelDelivery.checked;
        if (isDelivery) {
            const currentHour = new Date().getHours();
            if (currentHour < 19 || currentHour >= 22) {
                alert('Hostel delivery is only available between 7 PM and 10 PM');
                return;
            }
            if (!this.roomNumber.value) {
                alert('Please enter your room number');
                return;
            }
        }

        try {
            const orderData = {
                items: this.items.map(item => ({
                    menuItemID: item.ID,
                    quantity: item.quantity
                })),
                isDelivery,
                roomNumber: isDelivery ? this.roomNumber.value : ''
            };

            const order = await api.createOrder(orderData);
            this.items = [];
            this.updateCartUI();
            this.hideCart();
            alert(`Order placed successfully! Your order number is ${order.OrderNumber}`);
            window.location.href = '#orders';
        } catch (error) {
            alert(error.message);
        }
    }
}