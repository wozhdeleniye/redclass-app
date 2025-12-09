import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Container,
  Title,
  SimpleGrid,
  Card,
  Text,
  Badge,
  Group,
} from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { projectsApi } from '../api/projects';

export const ProjectsList: React.FC = () => {
  const navigate = useNavigate();

  // Получение моих проектов
  const { data: projects } = useQuery({
    queryKey: ['projects', 'my'],
    queryFn: () => projectsApi.getMy(),
  });

  return (
    <Container size="xl">
      <Title order={1} mb="lg">
        Мои проекты
      </Title>

      {projects && projects.length > 0 ? (
        <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
          {projects.map((project) => (
            <Card
              key={project.id}
              shadow="sm"
              padding="lg"
              radius="md"
              withBorder
              style={{ cursor: 'pointer' }}
              onClick={() => navigate(`/projects/${project.id}`)}
            >
              <Group justify="space-between" mb="xs">
                <Text fw={500}>{project.title}</Text>
                <Badge color="blue">{project.code}</Badge>
              </Group>
              <Text size="sm" c="dimmed" lineClamp={2}>
                {project.description}
              </Text>
            </Card>
          ))}
        </SimpleGrid>
      ) : (
        <Text c="dimmed" ta="center" py="xl">
          У вас пока нет проектов
        </Text>
      )}
    </Container>
  );
};

