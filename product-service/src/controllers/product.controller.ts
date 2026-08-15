import { Request, Response } from 'express';
import { Types } from 'mongoose';

import { Product } from '../models/product';
import { Category } from '../models/category';

// GET /api/products
export const getProducts = async (
  _req: Request,
  res: Response
): Promise<void> => {
  try {
    const products = await Product.find()
      .populate('categoryId')
      .sort({ name: 1 });

    res.status(200).json({
      products,
    });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error retrieving products',
    });
  }
};

// GET /api/products/:id
export const getProductById = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    // Validate product ID
    if (typeof id !== 'string' || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid product ID',
      });
      return;
    }

    const product = await Product.findById(id)
      .populate('categoryId');

    if (!product) {
      res.status(404).json({
        message: 'Product not found',
      });
      return;
    }

    res.status(200).json(product);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error retrieving product',
    });
  }
};

// POST /api/products
export const createProduct = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const {
      name,
      price,
      stock,
      categoryId,
      attributes,
    } = req.body;

    // Validate required fields
    if (!name || typeof name !== 'string' || name.trim() === '') {
      res.status(400).json({
        message: 'Product name is required',
      });
      return;
    }

    if (typeof price !== 'number' || price <= 0) {
      res.status(400).json({
        message: 'Product price is required and must be greater than 0',
      });
      return;
    }

    if (typeof stock !== 'number' || stock < 0) {
      res.status(400).json({
        message: 'Product stock is required and must be >= 0',
      });
      return;
    }

    // Validate categoryId
    if (
      typeof categoryId !== 'string' ||
      !Types.ObjectId.isValid(categoryId)
    ) {
      res.status(400).json({
        message: 'Invalid or missing category ID',
      });
      return;
    }

    // Check if category exists
    const category = await Category.findById(categoryId);

    if (!category) {
      res.status(404).json({
        message: 'Category not found',
      });
      return;
    }

    const product = await Product.create({
      name: name.trim(),
      price,
      stock,
      categoryId,
      attributes: attributes || {},
    });

    const productWithCategory = await product.populate('categoryId');

    res.status(201).json(productWithCategory);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error creating product',
    });
  }
};

// PUT /api/products/:id
export const updateProduct = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    // Validate product ID
    if (typeof id !== 'string' || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid product ID',
      });
      return;
    }

    const {
      name,
      price,
      stock,
      categoryId,
      attributes,
    } = req.body;

    // Validate categoryId if provided
    if (categoryId !== undefined) {
      if (
        typeof categoryId !== 'string' ||
        !Types.ObjectId.isValid(categoryId)
      ) {
        res.status(400).json({
          message: 'Invalid category ID',
        });
        return;
      }

      const category = await Category.findById(categoryId);

      if (!category) {
        res.status(404).json({
          message: 'Category not found',
        });
        return;
      }
    }

    console.log(`[MongoDB] PUT /api/products/${id} -> Operation: findByIdAndUpdate("${id}", {...})`);
    const product = await Product.findByIdAndUpdate(
      id,
      {
        name,
        price,
        stock,
        categoryId,
        attributes,
      },
      {
        new: true,
        runValidators: true,
      }
    ).populate('categoryId');

    if (!product) {
      console.log(`[MongoDB] ✗ Product not found for update`);
      res.status(404).json({
        message: 'Product not found',
      });
      return;
    }

    console.log(`[MongoDB] ✓ Product updated`);
    res.status(200).json(product);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error updating product',
    });
  }
};

// PATCH /api/products/:id/stock
export const updateProductStock = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;
    const { stock } = req.body;

    // Validate product ID
    if (typeof id !== 'string' || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid product ID',
      });
      return;
    }

    // Validate stock
    if (
      typeof stock !== 'number' ||
      !Number.isInteger(stock) ||
      stock < 0
    ) {
      res.status(400).json({
        message: 'Stock must be a non-negative integer',
      });
      return;
    }

    console.log(`[MongoDB] PATCH /api/products/${id}/stock -> Operation: findByIdAndUpdate("${id}", {stock: ${stock}})`)
    const product = await Product.findByIdAndUpdate(
      id,
      {
        stock,
      },
      {
        new: true,
        runValidators: true,
      }
    ).populate('categoryId');

    if (!product) {
      console.log(`[MongoDB] ✗ Product not found for stock update`);
      res.status(404).json({
        message: 'Product not found',
      });
      return;
    }

    console.log(`[MongoDB] ✓ Product stock updated to ${stock}`);
    res.status(200).json(product);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error updating product stock',
    });
  }
};

// DELETE /api/products/:id
export const deleteProduct = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    // Validate product ID
    if (typeof id !== 'string' || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid product ID',
      });
      return;
    }

    console.log(`[MongoDB] DELETE /api/products/${id} -> Operation: findByIdAndDelete("${id}")`);
    const product = await Product.findByIdAndDelete(id);

    if (!product) {
      console.log(`[MongoDB] ✗ Product not found for deletion`);
      res.status(404).json({
        message: 'Product not found',
      });
      return;
    }

    console.log(`[MongoDB] ✓ Product deleted`);
    res.status(200).json({
      message: 'Product deleted successfully',
      product,
    });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error deleting product',
    });
  }
};