import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Container,
  Title,
  Text,
  Button,
  Group,
  Modal,
  TextInput,
  Textarea,
  Alert,
  Paper,
  Stack,
  Card,
  Badge,
  Switch,
  MultiSelect,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { problemsApi } from '../api/problems';
import { projectsApi } from '../api/projects';
import {
  createProblemSchema,
  updateProblemSchema,
  type CreateProblemFormData,
  type UpdateProblemFormData,
} from '../schemas/problem';
import {
  createResultSchema,
  type CreateResultFormData,
} from '../schemas/result';
import dayjs from 'dayjs';

export const ProblemDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [createOpened, { open: openCreate, close: closeCreate }] = useDisclosure(false);
  const [editOpened, { open: openEdit, close: closeEdit }] = useDisclosure(false);
  const [resultOpened, { open: openResult, close: closeResult }] = useDisclosure(false);

  // Получение проблемы
  const { data: problemData } = useQuery({
    queryKey: ['problems', id],
    queryFn: () => problemsApi.getById(id!),
    enabled: !!id,
  });

  const problem = problemData?.problem;
  const result = problemData?.result;
  const childrenStats = problemData?.children_statistics;

  // Получение пользователей проекта
  const { data: projectUsers } = useQuery({
    queryKey: ['projects', problem?.project_id, 'users'],
    queryFn: () => projectsApi.getUsers(problem!.project_id),
    enabled: !!problem?.project_id,
  });

  // Получение подпроблем
  const { data: subproblems, isLoading: isLoadingSubproblems } = useQuery({
    queryKey: ['problems', id, 'subproblems'],
    queryFn: () => problemsApi.getSubproblems(id!),
    enabled: !!id,
  });

  // Создание подпроблемы
  const createSubproblemMutation = useMutation({
    mutationFn: (data: CreateProblemFormData) =>
      problemsApi.createSubproblem(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['problems', id] });
      queryClient.invalidateQueries({ queryKey: ['problems', id, 'subproblems'] });
      closeCreate();
      setCreateForm({
        title: '',
        description: '',
        start_time: '',
        end_time: '',
        assignee_ids: [],
      });
      setCreateErrors({});
    },
    onError: (error: any) => {
      setCreateApiError(error.response?.data?.message || 'Ошибка при создании подпроблемы');
    },
  });

  // Обновление проблемы
  const updateMutation = useMutation({
    mutationFn: (data: UpdateProblemFormData) => problemsApi.update(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['problems', id] });
      queryClient.invalidateQueries({ queryKey: ['problems', id, 'subproblems'] });
      closeEdit();
    },
    onError: (error: any) => {
      setEditApiError(error.response?.data?.message || 'Ошибка при обновлении');
    },
  });

  // Создание результата
  const createResultMutation = useMutation({
    mutationFn: (data: CreateResultFormData) =>
      problemsApi.createResult(id!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['problems', id] });
      closeResult();
      setResultForm({ done: false, comment: '' });
      setResultErrors({});
    },
    onError: (error: any) => {
      setResultApiError(error.response?.data?.message || 'Ошибка при создании результата');
    },
  });

  const [createForm, setCreateForm] = useState<CreateProblemFormData>({
    title: '',
    description: '',
    start_time: '',
    end_time: '',
    assignee_ids: [],
  });
  const [createErrors, setCreateErrors] = useState<Partial<Record<keyof CreateProblemFormData, string>>>({});
  const [createApiError, setCreateApiError] = useState('');

  const [editForm, setEditForm] = useState<UpdateProblemFormData>({});
  const [editErrors, setEditErrors] = useState<Partial<Record<keyof UpdateProblemFormData, string>>>({});
  const [editApiError, setEditApiError] = useState('');

  const [resultForm, setResultForm] = useState<CreateResultFormData>({
    done: false,
    comment: '',
  });
  const [resultErrors, setResultErrors] = useState<Partial<Record<keyof CreateResultFormData, string>>>({});
  const [resultApiError, setResultApiError] = useState('');

  const handleCreateSubmit = () => {
    setCreateErrors({});
    setCreateApiError('');

    const result = createProblemSchema.safeParse(createForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof CreateProblemFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof CreateProblemFormData] = issue.message;
        }
      });
      setCreateErrors(fieldErrors);
      return;
    }

    // Преобразуем даты в ISO формат с часовым поясом
    const formattedData = {
      ...createForm,
      start_time: dayjs(createForm.start_time).toISOString(),
      end_time: dayjs(createForm.end_time).toISOString(),
    };

    createSubproblemMutation.mutate(formattedData);
  };

  const handleEditClick = () => {
    if (problem) {
      setEditForm({
        title: problem.title,
        description: problem.description,
        start_time: dayjs(problem.start_time).format('YYYY-MM-DDTHH:mm'),
        end_time: dayjs(problem.end_time).format('YYYY-MM-DDTHH:mm'),
        assignee_ids: problem.assignees?.map((assignee) => assignee.user_id) || [],
      });
      openEdit();
    }
  };

  const handleEditSubmit = () => {
    setEditErrors({});
    setEditApiError('');

    const result = updateProblemSchema.safeParse(editForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof UpdateProblemFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof UpdateProblemFormData] = issue.message;
        }
      });
      setEditErrors(fieldErrors);
      return;
    }

    // Преобразуем даты в ISO формат с часовым поясом (только если они указаны)
    const formattedData = {
      ...editForm,
      ...(editForm.start_time && { start_time: dayjs(editForm.start_time).toISOString() }),
      ...(editForm.end_time && { end_time: dayjs(editForm.end_time).toISOString() }),
    };

    updateMutation.mutate(formattedData);
  };

  const handleResultSubmit = () => {
    setResultErrors({});
    setResultApiError('');

    const validationResult = createResultSchema.safeParse(resultForm);
    if (!validationResult.success) {
      const fieldErrors: Partial<Record<keyof CreateResultFormData, string>> = {};
      validationResult.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof CreateResultFormData] = issue.message;
        }
      });
      setResultErrors(fieldErrors);
      return;
    }

    createResultMutation.mutate(resultForm);
  };

  if (!problem) {
    return (
      <Container>
        <Text>Загрузка...</Text>
      </Container>
    );
  }

  const isRootProblem = !problem.parent_id;
  const hasSubproblems = subproblems && subproblems.length > 0;

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <div>
          <Group gap="xs" mb="xs">
            <Title order={1}>{problem.title}</Title>
            <Badge 
              size="lg" 
              variant="filled" 
              color={isRootProblem ? 'blue' : 'green'}
            >
              {isRootProblem ? 'Родительская проблема' : 'Дочерняя проблема'}
            </Badge>
          </Group>
          <Text c="dimmed" mt="xs">
            {problem.description}
          </Text>
        </div>
        <Group>
          <Button variant="light" onClick={handleEditClick}>Редактировать</Button>
          {!result && (
            <Button variant="outline" color="green" onClick={openResult}>
              Создать результат
            </Button>
          )}
          <Button onClick={openCreate}>Создать подпроблему</Button>
        </Group>
      </Group>

      {/* Навигация по иерархии */}
      {problem.parent_id && (
        <Paper p="md" withBorder mb="lg" bg="blue.0">
          <Group justify="space-between">
            <div>
              <Text fw={500} size="sm" c="blue.7" mb="xs">
                ⬆️ Родительская проблема
              </Text>
              <Text size="sm" c="dimmed">
                Эта проблема является подпроблемой
              </Text>
            </div>
            <Button
              variant="light"
              color="blue"
              onClick={() => navigate(`/problems/${problem.parent_id}`)}
            >
              Перейти к родителю
            </Button>
          </Group>
        </Paper>
      )}

      <Paper p="md" withBorder mb="lg">
        <Stack gap="xs">
          <Group>
            <Text fw={500}>Начало:</Text>
            <Text>{dayjs(problem.start_time).format('DD.MM.YYYY HH:mm')}</Text>
          </Group>
          <Group>
            <Text fw={500}>Окончание:</Text>
            <Text>{dayjs(problem.end_time).format('DD.MM.YYYY HH:mm')}</Text>
          </Group>
          <Group>
            <Text fw={500}>Статус:</Text>
            <Badge color={problem.solved ? 'green' : 'gray'}>
              {problem.solved ? 'Решена' : 'В работе'}
            </Badge>
          </Group>
          <div>
            <Text fw={500} mb="xs">
              👥 Исполнители:
            </Text>
            {problem.assignees && problem.assignees.length > 0 ? (
              <Group gap="xs">
                {problem.assignees.map((assignee) => (
                  <Badge key={assignee.id} variant="light" color="blue" size="lg">
                    {assignee.user.nickname}
                  </Badge>
                ))}
              </Group>
            ) : (
              <Text size="sm" c="dimmed">
                Не назначены
              </Text>
            )}
          </div>
        </Stack>
      </Paper>

      {/* Результат проблемы */}
      {result && (
        <Paper p="md" withBorder mb="lg" bg={result.done ? 'green.0' : 'yellow.0'}>
          <Title order={3} mb="md">
            Результат
          </Title>
          <Stack gap="xs">
            <Group>
              <Text fw={500}>Завершено:</Text>
              <Badge color={result.done ? 'green' : 'yellow'}>
                {result.done ? 'Да' : 'Нет'}
              </Badge>
            </Group>
            <div>
              <Text fw={500} mb="xs">
                Комментарий:
              </Text>
              <Text size="sm">{result.comment}</Text>
            </div>
            <Text size="xs" c="dimmed">
              Создано: {dayjs(result.created_at).format('DD.MM.YYYY HH:mm')}
            </Text>
          </Stack>
        </Paper>
      )}

      {/* Подпроблемы */}
      <Paper p="md" withBorder mb="lg">
        <Group justify="space-between" mb="md">
          <Title order={2}>
            Подпроблемы {hasSubproblems && `(${subproblems.length})`}
          </Title>
          {childrenStats && childrenStats.total > 0 && (
            <Group gap="xs">
              <Badge color="green" variant="filled" size="lg">
                ✓ {childrenStats.completed}
              </Badge>
              <Badge color="gray" variant="filled" size="lg">
                ⏳ {childrenStats.incomplete}
              </Badge>
              <Badge color="blue" variant="light" size="lg">
                Всего: {childrenStats.total}
              </Badge>
            </Group>
          )}
        </Group>

        {isLoadingSubproblems ? (
          <Text c="dimmed" ta="center" py="md">
            Загрузка подпроблем...
          </Text>
        ) : hasSubproblems ? (
          <Stack gap="md">
            {subproblems.map((subproblem) => (
              <Card
                key={subproblem.id}
                shadow="sm"
                padding="lg"
                radius="md"
                withBorder
                style={{ cursor: 'pointer', transition: 'transform 0.2s' }}
                onClick={() => navigate(`/problems/${subproblem.id}`)}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'translateY(-2px)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                }}
              >
                <Group justify="space-between" mb="sm">
                  <Group gap="xs">
                    <Badge variant="light" color="blue" size="lg">
                      #{subproblem.number}
                    </Badge>
                    <Badge
                      color={subproblem.solved ? 'green' : 'gray'}
                      variant="filled"
                    >
                      {subproblem.solved ? '✓ Решена' : '⏳ В работе'}
                    </Badge>
                  </Group>
                  <Button
                    variant="subtle"
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/problems/${subproblem.id}`);
                    }}
                  >
                    Открыть →
                  </Button>
                </Group>

                <Text fw={600} size="lg" mb="xs">
                  {subproblem.title}
                </Text>

                <Text size="sm" c="dimmed" mb="md" lineClamp={2}>
                  {subproblem.description}
                </Text>

                <Group gap="xl">
                  <div>
                    <Text size="xs" c="dimmed" mb={4}>
                      📅 Период
                    </Text>
                    <Text size="sm" fw={500}>
                      {dayjs(subproblem.start_time).format('DD.MM.YYYY')} -{' '}
                      {dayjs(subproblem.end_time).format('DD.MM.YYYY')}
                    </Text>
                  </div>

                  {subproblem.assignees && subproblem.assignees.length > 0 && (
                    <div>
                      <Text size="xs" c="dimmed" mb={4}>
                        👥 Исполнители
                      </Text>
                      <Group gap="xs">
                        {subproblem.assignees.slice(0, 3).map((assignee) => (
                          <Badge key={assignee.id} variant="outline" size="sm">
                            {assignee.user.nickname}
                          </Badge>
                        ))}
                        {subproblem.assignees.length > 3 && (
                          <Badge variant="outline" size="sm" color="gray">
                            +{subproblem.assignees.length - 3}
                          </Badge>
                        )}
                      </Group>
                    </div>
                  )}
                </Group>
              </Card>
            ))}
          </Stack>
        ) : (
          <Text c="dimmed" ta="center" py="xl">
            Подпроблем пока нет. Создайте первую подпроблему, чтобы разбить задачу на части.
          </Text>
        )}
      </Paper>

      {/* Модальное окно создания подпроблемы */}
      <Modal opened={createOpened} onClose={closeCreate} title="Создать подпроблему" size="lg">
        {createApiError && (
          <Alert color="red" mb="md">
            {createApiError}
          </Alert>
        )}

        <TextInput
          label="Название"
          placeholder="Подпроблема: Настроить Google OAuth"
          value={createForm.title}
          onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
          error={createErrors.title}
          mb="md"
        />

        <Textarea
          label="Описание"
          placeholder="Получить API ключи и настроить интеграцию"
          value={createForm.description}
          onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
          error={createErrors.description}
          mb="md"
          minRows={3}
        />

        <TextInput
          label="Дата начала"
          type="datetime-local"
          value={createForm.start_time}
          onChange={(e) => setCreateForm({ ...createForm, start_time: e.target.value })}
          error={createErrors.start_time}
          mb="md"
        />

        <TextInput
          label="Дата окончания"
          type="datetime-local"
          value={createForm.end_time}
          onChange={(e) => setCreateForm({ ...createForm, end_time: e.target.value })}
          error={createErrors.end_time}
          mb="md"
        />

        <MultiSelect
          label="Исполнители"
          placeholder="Выберите исполнителей"
          data={
            projectUsers?.map((user) => ({
              value: user.id,
              label: user.nickname,
            })) || []
          }
          value={createForm.assignee_ids}
          onChange={(value) =>
            setCreateForm({ ...createForm, assignee_ids: value })
          }
          searchable
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeCreate}>
            Отмена
          </Button>
          <Button onClick={handleCreateSubmit} loading={createSubproblemMutation.isPending}>
            Создать
          </Button>
        </Group>
      </Modal>

      {/* Модальное окно редактирования */}
      <Modal opened={editOpened} onClose={closeEdit} title="Редактировать проблему" size="lg">
        {editApiError && (
          <Alert color="red" mb="md">
            {editApiError}
          </Alert>
        )}

        <TextInput
          label="Название"
          value={editForm.title || ''}
          onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
          error={editErrors.title}
          mb="md"
        />

        <Textarea
          label="Описание"
          value={editForm.description || ''}
          onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
          error={editErrors.description}
          mb="md"
          minRows={3}
        />

        <TextInput
          label="Дата начала"
          type="datetime-local"
          value={editForm.start_time || ''}
          onChange={(e) => setEditForm({ ...editForm, start_time: e.target.value })}
          error={editErrors.start_time}
          mb="md"
        />

        <TextInput
          label="Дата окончания"
          type="datetime-local"
          value={editForm.end_time || ''}
          onChange={(e) => setEditForm({ ...editForm, end_time: e.target.value })}
          error={editErrors.end_time}
          mb="md"
        />

        <MultiSelect
          label="Исполнители"
          placeholder="Выберите исполнителей"
          data={
            projectUsers?.map((user) => ({
              value: user.id,
              label: user.nickname,
            })) || []
          }
          value={editForm.assignee_ids || []}
          onChange={(value) =>
            setEditForm({ ...editForm, assignee_ids: value })
          }
          searchable
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeEdit}>
            Отмена
          </Button>
          <Button onClick={handleEditSubmit} loading={updateMutation.isPending}>
            Сохранить
          </Button>
        </Group>
      </Modal>

      {/* Модальное окно создания результата */}
      <Modal opened={resultOpened} onClose={closeResult} title="Создать результат" size="lg">
        {resultApiError && (
          <Alert color="red" mb="md">
            {resultApiError}
          </Alert>
        )}

        <Switch
          label="Завершено"
          checked={resultForm.done}
          onChange={(e) => setResultForm({ ...resultForm, done: e.currentTarget.checked })}
          mb="md"
        />

        <Textarea
          label="Комментарий"
          placeholder="Описание результата работы..."
          value={resultForm.comment}
          onChange={(e) => setResultForm({ ...resultForm, comment: e.target.value })}
          error={resultErrors.comment}
          mb="md"
          minRows={4}
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeResult}>
            Отмена
          </Button>
          <Button onClick={handleResultSubmit} loading={createResultMutation.isPending} color="green">
            Создать
          </Button>
        </Group>
      </Modal>
    </Container>
  );
};

