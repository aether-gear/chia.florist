import type { Product } from '../models/Product'

export const productService = {
  async getProducts(): Promise<Product[]> {
    // Simulated API call
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          {
            id: '1',
            name: 'Rose Bouquet',
            description: 'A beautiful bouquet of red roses.',
            price: 49.99,
            imageUrl: 'https://placehold.co/400x400?text=Rose',
            category: 'flower-boards'
          },
          {
            id: '2',
            name: 'Wedding Board',
            description: 'Custom flower board for weddings.',
            price: 159.99,
            imageUrl: 'https://placehold.co/400x400?text=Wedding',
            category: 'flower-boards'
          }
        ])
      }, 500)
    })
  },

  async getProductById(id: string): Promise<Product | null> {
    // Simulated API call
    const products = await this.getProducts()
    return products.find(p => p.id === id) || null
  }
}
