import React from 'react';
import { Container, Title, Text, Paper, SimpleGrid, Button } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();

  return (
    <Container size="lg" py="xl">
      <Title order={1} mb="md">
        Добро пожаловать, {user?.nickname}!
      </Title>
      <Text c="dimmed" mb="xl">
        Выберите раздел для работы
      </Text>

      <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="lg">
        <Paper shadow="sm" p="xl" radius="md" withBorder>
          <Title order={3} mb="md">
            Предметы
          </Title>
          <Text c="dimmed" mb="md">
            Просматривайте, создавайте и управляйте предметами
          </Text>
          <Button onClick={() => navigate('/subjects')}>
            Перейти к предметам
          </Button>
        </Paper>

        <Paper shadow="sm" p="xl" radius="md" withBorder>
          <Title order={3} mb="md">
            Проекты
          </Title>
          <Text c="dimmed" mb="md">
            Управляйте своими проектами и задачами
          </Text>
          <Button onClick={() => navigate('/projects')}>
            Перейти к проектам
          </Button>
        </Paper>
      </SimpleGrid>
    </Container>
  );
};

