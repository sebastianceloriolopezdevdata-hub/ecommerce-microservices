import mongoose from "mongoose";

export const connectMongoDB = async (): Promise<void> => {
  try {
    await mongoose.connect(process.env.MONGODB_CNN as string);

    console.log("MongoDB database is online...");
  } catch (error) {
    console.error(error);

    throw new Error("Error starting MongoDB database...");
  }
};