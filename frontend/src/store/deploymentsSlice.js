import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const fetchDeployments = createAsyncThunk(
  'deployments/fetchDeployments',
  async () => {
    const response = await axios.get(`${API_URL}/deployments`);
    return response.data;
  }
);

export const createDeployment = createAsyncThunk(
  'deployments/createDeployment',
  async (deployment) => {
    const response = await axios.post(`${API_URL}/deployments`, deployment);
    return response.data;
  }
);

export const updateDeployment = createAsyncThunk(
  'deployments/updateDeployment',
  async ({ id, deployment }) => {
    const response = await axios.put(`${API_URL}/deployments/${id}`, deployment);
    return response.data;
  }
);

const deploymentsSlice = createSlice({
  name: 'deployments',
  initialState: {
    items: [],
    status: 'idle',
    error: null,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeployments.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(fetchDeployments.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.items = action.payload;
      })
      .addCase(fetchDeployments.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.error.message;
      })
      .addCase(createDeployment.fulfilled, (state, action) => {
        state.items.push(action.payload);
      })
      .addCase(updateDeployment.fulfilled, (state, action) => {
        const index = state.items.findIndex(d => d.id === action.payload.id);
        if (index !== -1) {
          state.items[index] = action.payload;
        }
      });
  },
});

export default deploymentsSlice.reducer; 