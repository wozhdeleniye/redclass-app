import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Container,
  Title,
  Text,
  Group,
  Paper,
  Switch,
} from '@mantine/core';
import { problemsApi } from '../api/problems';
import { ProblemTree } from '../components/ProblemTree';
import { ProjectStatistics } from '../components/ProjectStatistics';

export const ProjectDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [assignedOnly, setAssignedOnly] = useState(false);

  // Получение проблем проекта
  const { data: problems } = useQuery({
    queryKey: ['problems', 'project', id, assignedOnly],
    queryFn: () => problemsApi.getByProject(id!, assignedOnly),
    enabled: !!id,
  });

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <div>
          <Title order={1}>Детали проекта</Title>
          <Text c="dimmed" mt="xs">
            Управление проблемами и задачами проекта
          </Text>
        </div>
      </Group>

      {/* Статистика проекта */}
      <ProjectStatistics projectId={id!} />

      <Group justify="space-between" mb="md" mt="xl">
        <Title order={2}>Проблемы</Title>
        <Switch
          label="Только назначенные мне"
          checked={assignedOnly}
          onChange={(e) => setAssignedOnly(e.currentTarget.checked)}
        />
      </Group>

      {problems && problems.length > 0 ? (
        <ProblemTree problems={problems} projectId={id!} />
      ) : (
        <Paper p="xl" withBorder>
          <Text c="dimmed" ta="center">
            Главная проблема создается автоматически при создании проекта.
          </Text>
        </Paper>
      )}
    </Container>
  );
};

