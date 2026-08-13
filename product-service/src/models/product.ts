import { Schema, model, Document, Types } from 'mongoose';

// Interface for Product document
export interface IProduct extends Document {
  name: string;
  price: number;
  stock: number;
  categoryId: Types.ObjectId;
  attributes?: Record<string, any>;
}

const ProductSchema = new Schema<IProduct>(
  {
    name: {
      type: String,
      required: [true, 'The product name is required'],
      trim: true,
    },
    price: {
      type: Number,
      required: [true, 'The product price is required'],
      min: [0, 'The price cannot be negative'],
    },
    stock: {
      type: Number,
      required: [true, 'The product stock is required'],
      min: [0, 'The stock cannot be negative'],
    },
    categoryId: {
      type: Schema.Types.ObjectId,
      ref: 'Category',
      required: [true, 'The category is required'],
    },
    attributes: {
      type: Schema.Types.Mixed,
      default: {},
    },
  },
  {
    collection: 'products',
  }
);

ProductSchema.methods.toJSON = function () {
  const { __v, ...data } = this.toObject();
  return data;
};

export const Product = model<IProduct>('Product', ProductSchema);