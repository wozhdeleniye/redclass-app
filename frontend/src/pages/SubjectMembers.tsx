import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Container,
  Title,
  Button,
  Group,
  Table,
  Badge,
  Select,
  ActionIcon,
  Text,
} from '@mantine/core';
import { subjectsApi } from '../api/subjects';
import { rolesApi } from '../api/roles';
import type { RoleType } from '../types';

export const SubjectMembers: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Получение участников
  const { data: members } = useQuery({
    queryKey: ['subjects', id, 'members'],
    queryFn: () => subjectsApi.getMembers(id!),
    enabled: !!id,
  });

  // Изменение роли
  const changeRoleMutation = useMutation({
    mutationFn: ({ roleId, roleType }: { roleId: string; roleType: RoleType }) =>
      rolesApi.changeRole(id!, roleId, { role_type: roleType }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subjects', id, 'members'] });
    },
  });

  // Удаление пользователя
  const removeUserMutation = useMutation({
    mutationFn: (roleId: string) => rolesApi.removeUser(id!, roleId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['subjects', id, 'members'] });
    },
  });

  const getRoleBadgeColor = (role: RoleType) => {
    switch (role) {
      case 'admin':
        return 'red';
      case 'teacher':
        return 'blue';
      case 'student':
        return 'gray';
      default:
        return 'gray';
    }
  };

  const getRoleLabel = (role: RoleType) => {
    switch (role) {
      case 'admin':
        return 'Администратор';
      case 'teacher':
        return 'Преподаватель';
      case 'student':
        return 'Студент';
      default:
        return role;
    }
  };

  return (
    <Container size="xl">
      <Group justify="space-between" mb="lg">
        <Title order={1}>Участники предмета</Title>
        <Button variant="light" onClick={() => navigate(`/subjects/${id}`)}>
          Назад к предмету
        </Button>
      </Group>

      {members && members.length > 0 ? (
        <Table>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Имя</Table.Th>
              <Table.Th>Email</Table.Th>
              <Table.Th>Роль</Table.Th>
              <Table.Th>Действия</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {members.map((role) => (
              <Table.Tr key={role.id}>
                <Table.Td>{role.user?.nickname}</Table.Td>
                <Table.Td>{role.user?.email}</Table.Td>
                <Table.Td>
                  <Badge color={getRoleBadgeColor(role.role_type)}>
                    {getRoleLabel(role.role_type)}
                  </Badge>
                </Table.Td>
                <Table.Td>
                  <Group gap="xs">
                    <Select
                      value={role.role_type}
                      onChange={(value) => {
                        if (value) {
                          changeRoleMutation.mutate({
                            roleId: role.id,
                            roleType: value as RoleType,
                          });
                        }
                      }}
                      data={[
                        { value: 'student', label: 'Студент' },
                        { value: 'teacher', label: 'Преподаватель' },
                        { value: 'admin', label: 'Администратор' },
                      ]}
                      size="xs"
                      style={{ width: 160 }}
                    />
                    <ActionIcon
                      color="red"
                      variant="light"
                      onClick={() => removeUserMutation.mutate(role.id)}
                    >
                      ×
                    </ActionIcon>
                  </Group>
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      ) : (
        <Text c="dimmed" ta="center" py="xl">
          Нет участников
        </Text>
      )}
    </Container>
  );
};

