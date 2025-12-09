import { createBrowserRouter, Navigate } from 'react-router-dom';
import { ProtectedRoute } from './ProtectedRoute';
import { AppLayout } from '../layouts/AppLayout';
import { Login } from '../pages/Login';
import { Register } from '../pages/Register';
import { Dashboard } from '../pages/Dashboard';
import { SubjectsList } from '../pages/SubjectsList';
import { SubjectDetail } from '../pages/SubjectDetail';
import { SubjectMembers } from '../pages/SubjectMembers';
import { TaskDetail } from '../pages/TaskDetail';
import { ProjectsList } from '../pages/ProjectsList';
import { ProjectDetail } from '../pages/ProjectDetail';
import { ProblemDetail } from '../pages/ProblemDetail';

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <AppLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'subjects',
        element: <SubjectsList />,
      },
      {
        path: 'subjects/:id',
        element: <SubjectDetail />,
      },
      {
        path: 'subjects/:id/members',
        element: <SubjectMembers />,
      },
      {
        path: 'tasks/:id',
        element: <TaskDetail />,
      },
      {
        path: 'projects',
        element: <ProjectsList />,
      },
      {
        path: 'projects/:id',
        element: <ProjectDetail />,
      },
      {
        path: 'problems/:id',
        element: <ProblemDetail />,
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
]);

