import express, { Application, Request, Response } from 'express';
import cors from 'cors';

import { connectMongoDB } from '../database/mongoConnection';

import productRouter from '../routes/product.routes';
import categoryRouter from '../routes/categoryRoutes';

class Server {
  private app: Application;
  private port: string | number;

  constructor() {
    this.app = express();
    this.port = process.env.PORT || 3000;

    this.middlewares();
    this.routes();
    this.connectDB();

    this.app.get('/', (_req: Request, res: Response) => {
      res.json({
        ok: true,
        message: 'Product Service API running',
      });
    });
  }

  async connectDB(): Promise<void> {
    try {
      await connectMongoDB();
      console.log('MongoDB connection successful');
    } catch (error) {
      console.error('MongoDB connection failed', error);
    }
  }

  routes(): void {
    this.app.use('/api/products', productRouter);
    this.app.use('/api/categories', categoryRouter);
  }

  middlewares(): void {
    this.app.use(cors());
    this.app.use(express.json());
  }

  listen(): void {
    this.app.listen(this.port, () => {
      console.log(`Product Service running on port ${this.port}`);
    });
  }
}

export default Server;