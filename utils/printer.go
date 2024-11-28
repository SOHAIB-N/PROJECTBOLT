package utils

import (
	
	"fmt"
	"net"
	"restaurant-system/models"
	"time"
)

type PrinterConfig struct {
	Address string // IP address and port of the printer
	Port    int
}

type Printer struct {
	config PrinterConfig
	conn   net.Conn
}

func NewPrinter(config PrinterConfig) *Printer {
	return &Printer{
		config: config,
	}
}

func (p *Printer) connect() error {
	address := fmt.Sprintf("%s:%d", p.config.Address, p.config.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to printer: %v", err)
	}
	p.conn = conn
	return nil
}

func (p *Printer) disconnect() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Printer) sendCommand(cmd []byte) error {
	if p.conn == nil {
		if err := p.connect(); err != nil {
			return err
		}
	}
	_, err := p.conn.Write(cmd)
	return err
}

func PrintReceipt(order models.Order) error {
	// Initialize printer with default configuration
	// In production, load this from environment variables or configuration file
	printer := NewPrinter(PrinterConfig{
		Address: "192.168.1.100", // Replace with actual printer IP
		Port:    9100,            // Common port for thermal printers
	})
	defer printer.disconnect()

	// Format receipt content
	receipt := formatReceipt(order)

	// For demonstration, we'll print to console
	// In production, send to actual printer
	if err := printer.printToThermalPrinter(receipt); err != nil {
		// Fallback to console printing if printer is not available
		fmt.Println("Printer error, falling back to console:")
		fmt.Println(receipt)
		return err
	}

	return nil
}

func formatReceipt(order models.Order) string {
	receipt := fmt.Sprintf(`
============================
    RESTAURANT ORDER
============================
Order #: %s
Date: %s
----------------------------
`, order.OrderNumber, time.Now().Format("2006-01-02 15:04:05"))

	for _, item := range order.Items {
		itemLine := fmt.Sprintf("%dx %s\n    $%.2f each = $%.2f\n",
			item.Quantity,
			item.MenuItem.Name,
			item.Price,
			item.TotalPrice)
		receipt += itemLine
	}

	receipt += fmt.Sprintf(`
----------------------------
Subtotal: $%.2f
Tax: $%.2f
Total: $%.2f
============================
`, order.TotalAmount*0.9, order.TotalAmount*0.1, order.TotalAmount)

	if order.IsDelivery {
		receipt += fmt.Sprintf(`
Delivery Information:
Room: %s
============================
`, order.RoomNumber)
	}

	receipt += `
Thank you for your order!
Please keep this receipt.
============================
`

	return receipt
}

func (p *Printer) printToThermalPrinter(content string) error {
	// ESC/POS commands for thermal printer
	// Initialize printer
	init := []byte{0x1B, 0x40}
	if err := p.sendCommand(init); err != nil {
		return err
	}

	// Set text size (normal)
	textSize := []byte{0x1D, 0x21, 0x00}
	if err := p.sendCommand(textSize); err != nil {
		return err
	}

	// Send content
	if err := p.sendCommand([]byte(content)); err != nil {
		return err
	}

	// Cut paper
	cut := []byte{0x1D, 0x56, 0x41, 0x10}
	if err := p.sendCommand(cut); err != nil {
		return err
	}

	return nil
}