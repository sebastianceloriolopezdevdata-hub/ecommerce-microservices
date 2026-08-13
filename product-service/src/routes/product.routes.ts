import { Router } from 'express';

import {
  getProducts,
  getProductById,
  createProduct,
  updateProduct,
  updateProductStock,
  deleteProduct,
} from '../controllers/product.controller';

const router = Router();

router.get('/', getProducts);

router.get('/:id', getProductById);

router.post('/', createProduct);

router.put('/:id', updateProduct);

router.patch('/:id/stock', updateProductStock);

router.delete('/:id', deleteProduct);

export default router;