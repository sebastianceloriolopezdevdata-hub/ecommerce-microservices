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

    const category = await Category.create({
      name,
      description,
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

    const category = await Category.findByIdAndUpdate(
      id,
      {
        name,
        description,
      },
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
      res.status(409).json({
        message: 'Cannot delete category because it has products',
      });
      return;
    }

    const category = await Category.findByIdAndDelete(id);

    if (!category) {
      res.status(404).json({
        message: 'Category not found',
      });
      return;
    }

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