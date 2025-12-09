import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  Paper,
  Title,
  Text,
  Group,
  Stack,
  RingProgress,
  Badge,
  Loader,
} from '@mantine/core';
import { projectsApi } from '../api/projects';

interface ProjectStatisticsProps {
  projectId: string;
}

export const ProjectStatistics: React.FC<ProjectStatisticsProps> = ({
  projectId,
}) => {
  const { data: statistics, isLoading } = useQuery({
    queryKey: ['projects', projectId, 'statistics'],
    queryFn: () => projectsApi.getStatistics(projectId),
    enabled: !!projectId,
  });

  if (isLoading) {
    return (
      <Paper p="xl" withBorder>
        <Group justify="center">
          <Loader size="md" />
          <Text c="dimmed">Загрузка статистики...</Text>
        </Group>
      </Paper>
    );
  }

  if (!statistics) {
    return null;
  }

  const getProgressColor = (percentage: number) => {
    if (percentage >= 75) return 'green';
    if (percentage >= 50) return 'blue';
    if (percentage >= 25) return 'yellow';
    return 'red';
  };

  return (
    <Paper p="xl" withBorder shadow="sm">
      <Title order={2} mb="xl">
        📊 Статистика проекта
      </Title>

      <Group align="flex-start" gap="xl">
        {/* Круговая диаграмма */}
        <div>
          <RingProgress
            size={180}
            thickness={20}
            roundCaps
            sections={[
              {
                value: statistics.percentage,
                color: getProgressColor(statistics.percentage),
              },
            ]}
            label={
              <div style={{ textAlign: 'center' }}>
                <Text size="xl" fw={700}>
                  {statistics.percentage}%
                </Text>
                <Text size="xs" c="dimmed">
                  выполнено
                </Text>
              </div>
            }
          />
        </div>

        {/* Детальная статистика */}
        <Stack gap="md" style={{ flex: 1 }}>
          <Paper p="md" withBorder bg="green.0">
            <Group justify="space-between">
              <div>
                <Text size="sm" c="dimmed" mb={4}>
                  ✓ Выполнено
                </Text>
                <Text size="xl" fw={700} c="green.7">
                  {statistics.completed}
                </Text>
              </div>
              <Badge color="green" variant="filled" size="xl">
                {statistics.total > 0
                  ? Math.round((statistics.completed / statistics.total) * 100)
                  : 0}
                %
              </Badge>
            </Group>
          </Paper>

          <Paper p="md" withBorder bg="gray.0">
            <Group justify="space-between">
              <div>
                <Text size="sm" c="dimmed" mb={4}>
                  ⏳ В работе
                </Text>
                <Text size="xl" fw={700} c="gray.7">
                  {statistics.incomplete}
                </Text>
              </div>
              <Badge color="gray" variant="filled" size="xl">
                {statistics.total > 0
                  ? Math.round((statistics.incomplete / statistics.total) * 100)
                  : 0}
                %
              </Badge>
            </Group>
          </Paper>

          <Paper p="md" withBorder bg="blue.0">
            <Group justify="space-between">
              <div>
                <Text size="sm" c="dimmed" mb={4}>
                  📋 Всего проблем
                </Text>
                <Text size="xl" fw={700} c="blue.7">
                  {statistics.total}
                </Text>
              </div>
            </Group>
          </Paper>
        </Stack>
      </Group>

      {/* Прогресс бар внизу */}
      {statistics.total > 0 && (
        <div style={{ marginTop: '1.5rem' }}>
          <Group justify="space-between" mb="xs">
            <Text size="sm" fw={500}>
              Прогресс выполнения
            </Text>
            <Text size="sm" c="dimmed">
              {statistics.completed} из {statistics.total}
            </Text>
          </Group>
          <div
            style={{
              width: '100%',
              height: '12px',
              backgroundColor: '#e9ecef',
              borderRadius: '6px',
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                width: `${statistics.percentage}%`,
                height: '100%',
                backgroundColor: getProgressColor(statistics.percentage),
                transition: 'width 0.3s ease',
              }}
            />
          </div>
        </div>
      )}
    </Paper>
  );
};

