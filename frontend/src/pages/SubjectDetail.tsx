import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Container,
  Title,
  Text,
  Button,
  Group,
  Card,
  Badge,
  Modal,
  TextInput,
  Textarea,
  Alert,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { subjectsApi } from '../api/subjects';
import { tasksApi } from '../api/tasks';
import { createTaskSchema, type CreateTaskFormData } from '../schemas/task';
import dayjs from 'dayjs';

export const SubjectDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [createOpened, { open: openCreate, close: closeCreate }] = useDisclosure(false);

  // Получение данных предмета
  const { data: subject, isLoading } = useQuery({
    queryKey: ['subjects', id],
    queryFn: () => subjectsApi.getById(id!),
    enabled: !!id,
  });

  // Получение заданий предмета
  const { data: tasks } = useQuery({
    queryKey: ['tasks', 'subject', id],
    queryFn: () => tasksApi.getBySubject(id!),
    enabled: !!id,
  });

  // Создание задания
  const createTaskMutation = useMutation({
    mutationFn: (data: CreateTaskFormData) => tasksApi.create(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks', 'subject', id] });
      closeCreate();
      setCreateForm({ title: '', description: '', due_date: '' });
      setCreateErrors({});
    },
    onError: (error: any) => {
      setCreateApiError(error.response?.data?.message || 'Ошибка при создании задания');
    },
  });

  const [createForm, setCreateForm] = useState<CreateTaskFormData>({
    title: '',
    description: '',
    due_date: '',
  });
  const [createErrors, setCreateErrors] = useState<Partial<Record<keyof CreateTaskFormData, string>>>({});
  const [createApiError, setCreateApiError] = useState('');

  const handleCreateSubmit = () => {
    setCreateErrors({});
    setCreateApiError('');

    const result = createTaskSchema.safeParse(createForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof CreateTaskFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof CreateTaskFormData] = issue.message;
        }
      });
      setCreateErrors(fieldErrors);
      return;
    }

    // Преобразуем дату в ISO формат с часовым поясом
    const formattedData = {
      ...createForm,
      due_date: dayjs(createForm.due_date).toISOString(),
    };

    createTaskMutation.mutate(formattedData);
  };

  if (isLoading) {
    return (
      <Container>
        <Text>Загрузка...</Text>
      </Container>
    );
  }

  if (!subject) {
    return (
      <Container>
        <Text>Предмет не найден</Text>
      </Container>
    );
  }

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <div>
          <Group>
            <Title order={1}>{subject.name}</Title>
            <Badge size="lg">{subject.code}</Badge>
          </Group>
          <Text c="dimmed" mt="xs">
            {subject.description}
          </Text>
        </div>
        <Group>
          <Button variant="light" onClick={() => navigate(`/subjects/${id}/members`)}>
            Участники
          </Button>
          <Button onClick={openCreate}>Создать задание</Button>
        </Group>
      </Group>

      <Title order={2} mb="md">
        Задания
      </Title>

      {tasks?.data && tasks.data.length > 0 ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          {tasks.data.map((task) => (
            <Card key={task.id} shadow="sm" padding="lg" radius="md" withBorder>
              <Group justify="space-between">
                <div style={{ flex: 1 }}>
                  <Text fw={500} size="lg" mb="xs">
                    {task.title}
                  </Text>
                  <Text size="sm" c="dimmed" mb="xs">
                    {task.description}
                  </Text>
                  <Text size="sm" c="dimmed">
                    Срок: {dayjs(task.due_date).format('DD.MM.YYYY HH:mm')}
                  </Text>
                </div>
                <Button variant="light" onClick={() => navigate(`/tasks/${task.id}`)}>
                  Подробнее
                </Button>
              </Group>
            </Card>
          ))}
        </div>
      ) : (
        <Text c="dimmed" ta="center" py="xl">
          Заданий пока нет
        </Text>
      )}

      {/* Модальное окно создания задания */}
      <Modal opened={createOpened} onClose={closeCreate} title="Создать задание" size="lg">
        {createApiError && (
          <Alert color="red" mb="md">
            {createApiError}
          </Alert>
        )}

        <TextInput
          label="Название"
          placeholder="Решить задачи 1-10"
          value={createForm.title}
          onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
          error={createErrors.title}
          mb="md"
        />

        <Textarea
          label="Описание"
          placeholder="Решить задачи из учебника на странице 45"
          value={createForm.description}
          onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
          error={createErrors.description}
          mb="md"
          minRows={4}
        />

        <TextInput
          label="Срок выполнения"
          type="datetime-local"
          value={createForm.due_date}
          onChange={(e) => setCreateForm({ ...createForm, due_date: e.target.value })}
          error={createErrors.due_date}
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeCreate}>
            Отмена
          </Button>
          <Button onClick={handleCreateSubmit} loading={createTaskMutation.isPending}>
            Создать
          </Button>
        </Group>
      </Modal>
    </Container>
  );
};

