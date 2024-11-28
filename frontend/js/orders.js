class Orders {
    constructor() {
        this.ordersContainer = document.getElementById('main-content');
    }

    async loadOrders() {
        try {
            const orders = await api.getUserOrders();
            this.renderOrders(orders);
        } catch (error) {
            console.error('Error loading orders:', error);
            this.ordersContainer.innerHTML = '<p class="error">Error loading orders. Please try again later.</p>';
        }
    }

    renderOrders(orders) {
        if (orders.length === 0) {
            this.ordersContainer.innerHTML = '<p>No orders found.</p>';
            return;
        }

        const ordersHTML = orders.map(order => `
            <div class="order-card">
                <div class="order-header">
                    <h3>Order #${order.OrderNumber}</h3>
                    <span class="order-status ${order.Status}">${order.Status}</span>
                </div>
                <div class="order-items">
                    ${order.Items.map(item => `
                        <div class="order-item">
                            <span>${item.Quantity}x ${item.MenuItem.Name}</span>
                            <span>$${(item.Price * item.Quantity).toFixed(2)}</span>
                        </div>
                    `).join('')}
                </div>
                <div class="order-footer">
                    <p>Total: $${order.TotalAmount.toFixed(2)}</p>
                    ${order.IsDelivery ? `<p>Delivery to Room: ${order.RoomNumber}</p>` : ''}
                    <p>Ordered on: ${new Date(order.CreatedAt).toLocaleString()}</p>
                </div>
            </div>
        `).join('');

        this.ordersContainer.innerHTML = `
            <div class="orders-container">
                <h2>My Orders</h2>
                ${ordersHTML}
            </div>
        `;
    }
}