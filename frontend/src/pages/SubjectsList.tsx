import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Container,
  Title,
  Tabs,
  SimpleGrid,
  Button,
  Modal,
  TextInput,
  Textarea,
  Group,
  Card,
  Text,
  Badge,
  Alert,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { useNavigate } from 'react-router-dom';
import { subjectsApi } from '../api/subjects';
import {
  createSubjectSchema,
  joinSubjectSchema,
  type CreateSubjectFormData,
  type JoinSubjectFormData,
} from '../schemas/subject';

export const SubjectsList: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [createOpened, { open: openCreate, close: closeCreate }] = useDisclosure(false);
  const [joinOpened, { open: openJoin, close: closeJoin }] = useDisclosure(false);

  // Получение всех предметов
  const { data: allSubjects } = useQuery({
    queryKey: ['subjects', 'all'],
    queryFn: () => subjectsApi.getAll({ limit: 100, offset: 0 }),
  });

  // Получение моих предметов
  const { data: mySubjects } = useQuery({
    queryKey: ['subjects', 'my'],
    queryFn: () => subjectsApi.getMy(),
  });

  // Создание предмета
  const createMutation = useMutation({
    mutationFn: subjectsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subjects'] });
      closeCreate();
      setCreateForm({ name: '', description: '', code: '' });
      setCreateErrors({});
    },
    onError: (error: any) => {
      setCreateApiError(error.response?.data?.message || 'Ошибка при создании предмета');
    },
  });

  // Присоединение к предмету
  const joinMutation = useMutation({
    mutationFn: subjectsApi.join,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subjects'] });
      closeJoin();
      setJoinForm({ code: '' });
      setJoinErrors({});
    },
    onError: (error: any) => {
      setJoinApiError(error.response?.data?.message || 'Ошибка при присоединении');
    },
  });

  // Форма создания
  const [createForm, setCreateForm] = useState<CreateSubjectFormData>({
    name: '',
    description: '',
    code: '',
  });
  const [createErrors, setCreateErrors] = useState<Partial<Record<keyof CreateSubjectFormData, string>>>({});
  const [createApiError, setCreateApiError] = useState('');

  // Форма присоединения
  const [joinForm, setJoinForm] = useState<JoinSubjectFormData>({ code: '' });
  const [joinErrors, setJoinErrors] = useState<Partial<Record<keyof JoinSubjectFormData, string>>>({});
  const [joinApiError, setJoinApiError] = useState('');

  const handleCreateSubmit = () => {
    setCreateErrors({});
    setCreateApiError('');

    const result = createSubjectSchema.safeParse(createForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof CreateSubjectFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof CreateSubjectFormData] = issue.message;
        }
      });
      setCreateErrors(fieldErrors);
      return;
    }

    createMutation.mutate(createForm);
  };

  const handleJoinSubmit = () => {
    setJoinErrors({});
    setJoinApiError('');

    const result = joinSubjectSchema.safeParse(joinForm);
    if (!result.success) {
      const fieldErrors: Partial<Record<keyof JoinSubjectFormData, string>> = {};
      result.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          fieldErrors[issue.path[0] as keyof JoinSubjectFormData] = issue.message;
        }
      });
      setJoinErrors(fieldErrors);
      return;
    }

    joinMutation.mutate(joinForm);
  };

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <Title order={1}>Предметы</Title>
        <Group>
          <Button onClick={openJoin}>Присоединиться по коду</Button>
          <Button onClick={openCreate}>Создать предмет</Button>
        </Group>
      </Group>

      <Tabs defaultValue="my">
        <Tabs.List>
          <Tabs.Tab value="my">Мои предметы</Tabs.Tab>
          <Tabs.Tab value="all">Все предметы</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="my" pt="xl">
          <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
            {mySubjects?.map((subject) => (
              <Card
                key={subject.id}
                shadow="sm"
                padding="lg"
                radius="md"
                withBorder
                style={{ cursor: 'pointer' }}
                onClick={() => navigate(`/subjects/${subject.id}`)}
              >
                <Group justify="space-between" mb="xs">
                  <Text fw={500}>{subject.name}</Text>
                  <Badge color="blue">{subject.code}</Badge>
                </Group>
                <Text size="sm" c="dimmed" lineClamp={2}>
                  {subject.description}
                </Text>
              </Card>
            ))}
          </SimpleGrid>
          {mySubjects?.length === 0 && (
            <Text c="dimmed" ta="center" py="xl">
              У вас пока нет предметов
            </Text>
          )}
        </Tabs.Panel>

        <Tabs.Panel value="all" pt="xl">
          <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }}>
            {allSubjects?.data?.map((subject) => (
              <Card
                key={subject.id}
                shadow="sm"
                padding="lg"
                radius="md"
                withBorder
                style={{ cursor: 'pointer' }}
                onClick={() => navigate(`/subjects/${subject.id}`)}
              >
                <Group justify="space-between" mb="xs">
                  <Text fw={500}>{subject.name}</Text>
                  <Badge color="gray">{subject.code}</Badge>
                </Group>
                <Text size="sm" c="dimmed" lineClamp={2}>
                  {subject.description}
                </Text>
              </Card>
            ))}
          </SimpleGrid>
          {allSubjects?.data?.length === 0 && (
            <Text c="dimmed" ta="center" py="xl">
              Нет доступных предметов
            </Text>
          )}
        </Tabs.Panel>
      </Tabs>

      {/* Модальное окно создания предмета */}
      <Modal opened={createOpened} onClose={closeCreate} title="Создать предмет">
        {createApiError && (
          <Alert color="red" mb="md">
            {createApiError}
          </Alert>
        )}

        <TextInput
          label="Название"
          placeholder="Математика"
          value={createForm.name}
          onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
          error={createErrors.name}
          mb="md"
        />

        <Textarea
          label="Описание"
          placeholder="Курс высшей математики"
          value={createForm.description}
          onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
          error={createErrors.description}
          mb="md"
          minRows={3}
        />

        <TextInput
          label="Код предмета"
          placeholder="MATH101"
          value={createForm.code}
          onChange={(e) => setCreateForm({ ...createForm, code: e.target.value.toUpperCase() })}
          error={createErrors.code}
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeCreate}>
            Отмена
          </Button>
          <Button onClick={handleCreateSubmit} loading={createMutation.isPending}>
            Создать
          </Button>
        </Group>
      </Modal>

      {/* Модальное окно присоединения */}
      <Modal opened={joinOpened} onClose={closeJoin} title="Присоединиться к предмету">
        {joinApiError && (
          <Alert color="red" mb="md">
            {joinApiError}
          </Alert>
        )}

        <TextInput
          label="Код предмета"
          placeholder="MATH101"
          value={joinForm.code}
          onChange={(e) => setJoinForm({ code: e.target.value.toUpperCase() })}
          error={joinErrors.code}
          mb="md"
        />

        <Group justify="flex-end" mt="md">
          <Button variant="subtle" onClick={closeJoin}>
            Отмена
          </Button>
          <Button onClick={handleJoinSubmit} loading={joinMutation.isPending}>
            Присоединиться
          </Button>
        </Group>
      </Modal>
    </Container>
  );
};

