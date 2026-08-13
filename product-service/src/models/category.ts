import { Schema, model, Document } from 'mongoose';

// Interface for Category document
export interface ICategory extends Document {
  name: string;
  description?: string;
}

const CategorySchema = new Schema<ICategory>(
  {
    name: {
      type: String,
      required: [true, 'The category name is required'],
      unique: true,
      trim: true,
    },
    description: {
      type: String,
      trim: true,
    },
  },
  {
    collection: 'categories',
  }
);

CategorySchema.methods.toJSON = function () {
  const { __v, ...data } = this.toObject();
  return data;
};

export const Category = model<ICategory>('Category', CategorySchema);