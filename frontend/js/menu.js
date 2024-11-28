class Menu {
    constructor() {
        this.menuContainer = document.getElementById('main-content');
        this.loadMenu();
    }

    async loadMenu() {
        try {
            const menuItems = await api.getMenu();
            this.renderMenu(menuItems);
        } catch (error) {
            console.error('Error loading menu:', error);
            this.menuContainer.innerHTML = '<p class="error">Error loading menu. Please try again later.</p>';
        }
    }

    renderMenu(items) {
        const menuHTML = `
            <div class="menu-grid">
                ${items.map(item => this.renderMenuItem(item)).join('')}
            </div>
        `;
        this.menuContainer.innerHTML = menuHTML;
    }

    renderMenuItem(item) {
        return `
            <div class="menu-item" data-id="${item.ID}">
                <img src="${item.ImageURL || 'images/default-dish.jpg'}" alt="${item.Name}">
                <div class="menu-item-content">
                    <h3 class="menu-item-title">${item.Name}</h3>
                    <p class="menu-item-description">${item.Description}</p>
                    <p class="menu-item-price">$${item.Price.toFixed(2)}</p>
                    <button onclick="cart.addItem(${JSON.stringify(item)})">Add to Cart</button>
                </div>
            </div>
        `;
    }
}