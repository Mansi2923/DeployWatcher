import { configureStore } from '@reduxjs/toolkit';
import deploymentsReducer from './deploymentsSlice';

export const store = configureStore({
  reducer: {
    deployments: deploymentsReducer,
  },
}); 