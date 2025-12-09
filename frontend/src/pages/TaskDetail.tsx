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
  Modal,
  TextInput,
  Textarea,
  Alert,
  Badge,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { tasksApi } from '../api/tasks';
import { projectsApi } from '../api/projects';
import {
  createProjectSchema,
  joinProjectSchema,
  type CreateProjectFormData,
  type JoinProjectFormData,
} from '../schemas/project';
import dayjs from 'dayjs';

export const TaskDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [createOpened, { open: openCreate, close: closeCreate }] = useDisclosure(false);
  const [joinOpened, { open: openJoin, close: closeJoin }] = useDisclosure(false);

  // Получение задания
  const { data: task } = useQuery({
    queryKey: ['tasks', id],
    queryFn: () => tasksApi.getById(id!),
    enabled: !!id,
  });

  // Получение проектов задания
  const { data: projects } = useQuery({
    queryKey: ['projects', 'task', id],
    queryFn: () => projectsApi.getByTask(id!),
    enabled: !!id,
  });

  // Создание проекта
  const createProjectMutation = useMutation({
    mutationFn: (data: CreateProjectFormData) => projectsApi.create(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', 'task', id] });
      closeCreate();
      setCreateForm({ title: '', description: '' });
      setCreateErrors({});
    },
    onError: (error: any) => {
      setCreateApiError(error.response?.data?.message || 'Ошибка при создании проекта');
    },
  });

  // Присоединение к проекту
  const joinProjectMutation = useMutation({
    mutationFn: projectsApi.join,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      closeJoin();
      setJoinForm({ code: '' });
      setJoinErrors({});
    },
    onError: (error: any) => {
      setJoinApiError(error.response?.data?.message || 'Ошибка при присоединении');
    },
  });

  const [createForm, setCreateForm] = useState<CreateProjectFormData>({
    title: '',
    description: '',
  });
  const [createErrors, setCreateErrors] = useState<Partial<Record<keyof CreateProjectFormData, string>>>({});
  const [createApiError, setCreateApiError] = useState('');

  const [joinForm, setJoinForm] = useState<JoinProjectFormData>({ code: '' });
  const [joinErrors, setJoinErrors] = useState<Partial<Record<keyof JoinProjectFormData, string>>>({});
  const [joinApiError, setJoinApiError] = useState('');

  const handleCreateSubmit = () => {
    setCreateErrors({});
    setCreateApiError('');

    const result = createProjectSchema.safeParse(createForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof CreateProjectFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof CreateProjectFormData] = issue.message;
        }
      });
      setCreateErrors(fieldErrors);
      return;
    }

    createProjectMutation.mutate(createForm);
  };

  const handleJoinSubmit = () => {
    setJoinErrors({});
    setJoinApiError('');

    const result = joinProjectSchema.safeParse(joinForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof JoinProjectFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof JoinProjectFormData] = issue.message;
        }
      });
      setJoinErrors(fieldErrors);
      return;
    }

    joinProjectMutation.mutate(joinForm);
  };

  if (!task) {
    return (
      <Container>
        <Text>Загрузка...</Text>
      </Container>
    );
  }

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <div>
          <Title order={1}>{task.title}</Title>
          <Text c="dimmed" mt="xs">
            {task.description}
          </Text>
          <Text size="sm" c="dimmed" mt="xs">
            Срок: {dayjs(task.due_date).format('DD.MM.YYYY HH:mm')}
          </Text>
        </div>
        <Group>
          <Button variant="light" onClick={openJoin}>
            Присоединиться к проекту
          </Button>
          <Button onClick={openCreate}>Создать проект</Button>
        </Group>
      </Group>

      <Title order={2} mb="md">
        Проекты
      </Title>

      {projects && projects.length > 0 ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
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
              <Group justify="space-between">
                <div>
                  <Group>
                    <Text fw={500} size="lg">
                      {project.title}
                    </Text>
                    <Badge>{project.code}</Badge>
                  </Group>
                  <Text size="sm" c="dimmed" mt="xs">
                    {project.description}
                  </Text>
                </div>
                <Button variant="light">Открыть</Button>
              </Group>
            </Card>
          ))}
        </div>
      ) : (
        <Text c="dimmed" ta="center" py="xl">
          Проектов пока нет
        </Text>
      )}

      {/* Модальное окно создания проекта */}
      <Modal opened={createOpened} onClose={closeCreate} title="Создать проект" size="lg">
        {createApiError && (
          <Alert color="red" mb="md">
            {createApiError}
          </Alert>
        )}

        <TextInput
          label="Название"
          placeholder="Проект: Веб-приложение"
          value={createForm.title}
          onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
          error={createErrors.title}
          mb="md"
        />

        <Textarea
          label="Описание"
          placeholder="Разработка веб-приложения для управления задачами"
          value={createForm.description}
          onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
          error={createErrors.description}
          mb="md"
          minRows={4}
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeCreate}>
            Отмена
          </Button>
          <Button onClick={handleCreateSubmit} loading={createProjectMutation.isPending}>
            Создать
          </Button>
        </Group>
      </Modal>

      {/* Модальное окно присоединения */}
      <Modal opened={joinOpened} onClose={closeJoin} title="Присоединиться к проекту">
        {joinApiError && (
          <Alert color="red" mb="md">
            {joinApiError}
          </Alert>
        )}

        <TextInput
          label="Код проекта"
          placeholder="a1b2c3d4"
          value={joinForm.code}
          onChange={(e) => setJoinForm({ code: e.target.value })}
          error={joinErrors.code}
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeJoin}>
            Отмена
          </Button>
          <Button onClick={handleJoinSubmit} loading={joinProjectMutation.isPending}>
            Присоединиться
          </Button>
        </Group>
      </Modal>
    </Container>
  );
};

