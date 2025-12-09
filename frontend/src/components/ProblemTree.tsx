import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Text, Group, Button, Badge } from '@mantine/core';
import type { Problem } from '../types';
import dayjs from 'dayjs';

interface ProblemTreeProps {
  problems: Problem[];
  projectId: string;
  level?: number;
}

export const ProblemTree: React.FC<ProblemTreeProps> = ({
  problems,
  projectId,
  level = 0,
}) => {
  const navigate = useNavigate();

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      {problems.map((problem) => (
        <div key={problem.id} style={{ marginLeft: `${level * 2}rem` }}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Group justify="space-between">
              <div style={{ flex: 1 }}>
                <Group gap="xs" mb="xs">
                  <Badge variant="light" color="blue">
                    #{problem.number}
                  </Badge>
                  <Badge color={problem.solved ? 'green' : 'gray'}>
                    {problem.solved ? 'Решена' : 'В работе'}
                  </Badge>
                </Group>
                <Text fw={500} size="lg" mb="xs">
                  {problem.title}
                </Text>
                <Text size="sm" c="dimmed" mb="xs">
                  {problem.description}
                </Text>
                <Text size="sm" c="dimmed" mb="xs">
                  📅 {dayjs(problem.start_time).format('DD.MM.YYYY')} -{' '}
                  {dayjs(problem.end_time).format('DD.MM.YYYY')}
                </Text>
                <Group gap="xs" mt="xs">
                  <Text size="sm" c="dimmed">
                    👥 Исполнители:
                  </Text>
                  {problem.assignees && problem.assignees.length > 0 ? (
                    problem.assignees.map((assignee) => (
                      <Badge key={assignee.id} variant="outline" size="sm" color="blue">
                        {assignee.user.nickname}
                      </Badge>
                    ))
                  ) : (
                    <Text size="sm" c="dimmed">
                      не назначены
                    </Text>
                  )}
                </Group>
              </div>
              <Button
                variant="light"
                onClick={() => navigate(`/problems/${problem.id}`)}
              >
                Открыть
              </Button>
            </Group>
          </Card>

          {problem.subproblems && problem.subproblems.length > 0 && (
            <ProblemTree
              problems={problem.subproblems}
              projectId={projectId}
              level={level + 1}
            />
          )}
        </div>
      ))}
    </div>
  );
};

