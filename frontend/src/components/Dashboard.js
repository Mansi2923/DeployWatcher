import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
  Chip,
  CircularProgress,
} from '@mui/material';
import { format } from 'date-fns';
import { fetchDeployments } from '../store/deploymentsSlice';

const statusColors = {
  queued: 'default',
  'in-progress': 'primary',
  successful: 'success',
  failed: 'error',
};

function DeploymentCard({ deployment }) {
  return (
    <Card sx={{ mb: 2 }}>
      <CardContent>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Box display="flex" justifyContent="space-between" alignItems="center">
              <Typography variant="h6">{deployment.appName}</Typography>
              <Chip
                label={deployment.status}
                color={statusColors[deployment.status]}
                size="small"
              />
            </Box>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="body2" color="text.secondary">
              Environment: {deployment.environment}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="body2" color="text.secondary">
              Branch: {deployment.branch}
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="body2" color="text.secondary">
              Started: {format(new Date(deployment.startedAt), 'PPpp')}
            </Typography>
          </Grid>
          {deployment.completedAt && (
            <Grid item xs={12}>
              <Typography variant="body2" color="text.secondary">
                Completed: {format(new Date(deployment.completedAt), 'PPpp')}
              </Typography>
            </Grid>
          )}
        </Grid>
      </CardContent>
    </Card>
  );
}

function Dashboard() {
  const dispatch = useDispatch();
  const { items: deployments, status, error } = useSelector((state) => state.deployments);

  useEffect(() => {
    dispatch(fetchDeployments());
    // Poll for updates every 30 seconds
    const interval = setInterval(() => {
      dispatch(fetchDeployments());
    }, 30000);
    return () => clearInterval(interval);
  }, [dispatch]);

  if (status === 'loading') {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (status === 'failed') {
    return (
      <Typography color="error" align="center">
        Error: {error}
      </Typography>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Deployment Dashboard
      </Typography>
      <Grid container spacing={2}>
        {deployments.map((deployment) => (
          <Grid item xs={12} md={6} lg={4} key={deployment.id}>
            <DeploymentCard deployment={deployment} />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
}

export default Dashboard; 