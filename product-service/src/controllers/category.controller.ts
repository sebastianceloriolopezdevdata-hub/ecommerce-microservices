import { Request, Response } from 'express';
import { Types } from 'mongoose';

import { Category } from '../models/category';
import { Product } from '../models/product';

// GET /api/categories
export const getCategories = async (
  _req: Request,
  res: Response
): Promise<void> => {
  try {
    const categories = await Category.find()
      .sort({ name: 1 });

    res.status(200).json({
      categories,
    });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error retrieving categories',
    });
  }
};

// GET /api/categories/:id
export const getCategoryById = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    if (!id || Array.isArray(id) || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid category ID',
      });
      return;
    }

    const category = await Category.findById(id);

    if (!category) {
      res.status(404).json({
        message: 'Category not found',
      });
      return;
    }

    res.status(200).json(category);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error retrieving category',
    });
  }
};

// POST /api/categories
export const createCategory = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { name, description } = req.body;

    if (!name || typeof name !== 'string' || name.trim() === '') {
      res.status(400).json({
        message: 'Category name is required',
      });
      return;
    }

    const category = await Category.create({
      name: name.trim(),
      description: description || '',
    });

    res.status(201).json(category);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error creating category',
    });
  }
};

// PUT /api/categories/:id
export const updateCategory = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;
    const { name, description } = req.body;

    if (!id || Array.isArray(id) || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid category ID',
      });
      return;
    }

    // Validate name if provided
    if (name !== undefined && (typeof name !== 'string' || name.trim() === '')) {
      res.status(400).json({
        message: 'Category name must be a non-empty string',
      });
      return;
    }

    const updateData: any = {};
    if (name !== undefined) updateData.name = name.trim();
    if (description !== undefined) updateData.description = description;

    const category = await Category.findByIdAndUpdate(
      id,
      updateData,
      {
        new: true,
        runValidators: true,
      }
    );

    if (!category) {
      res.status(404).json({
        message: 'Category not found',
      });
      return;
    }

    res.status(200).json(category);
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error updating category',
    });
  }
};

// DELETE /api/categories/:id
export const deleteCategory = async (
  req: Request,
  res: Response
): Promise<void> => {
  try {
    const { id } = req.params;

    if (!id || Array.isArray(id) || !Types.ObjectId.isValid(id)) {
      res.status(400).json({
        message: 'Invalid category ID',
      });
      return;
    }

    // Check if category has products
    const productsUsingCategory = await Product.countDocuments({
      categoryId: id,
    });

    if (productsUsingCategory > 0) {
      console.log(`[MongoDB] DELETE /api/categories/${id} -> Cannot delete (${productsUsingCategory} products found)`);
      res.status(409).json({
        message: 'Cannot delete category because it has products',
      });
      return;
    }

    console.log(`[MongoDB] DELETE /api/categories/${id} -> Operation: findByIdAndDelete("${id}")`);
    const category = await Category.findByIdAndDelete(id);

    if (!category) {
      console.log(`[MongoDB] ✗ Category not found for deletion`);
      res.status(404).json({
        message: 'Category not found',
      });
      return;
    }

    console.log(`[MongoDB] ✓ Category deleted`);
    res.status(200).json({
      message: 'Category deleted successfully',
      category,
    });
  } catch (error) {
    console.error(error);

    res.status(500).json({
      message: 'Error deleting category',
    });
  }
};