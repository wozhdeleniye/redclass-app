import React from 'react';
import { Outlet, useNavigate, Link } from 'react-router-dom';
import {
  AppShell,
  Burger,
  Group,
  Button,
  NavLink,
  Text,
  Avatar,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { useAuth } from '../contexts/AuthContext';

export const AppLayout: React.FC = () => {
  const [opened, { toggle }] = useDisclosure();
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: 'sm',
        collapsed: { mobile: !opened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Group>
            <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
            <Text size="xl" fw={700}>
              RedClass
            </Text>
          </Group>

          <Group>
            <Group gap="xs">
              <Avatar color="blue" radius="xl">
                {user?.nickname.charAt(0).toUpperCase()}
              </Avatar>
              <Text size="sm">{user?.nickname}</Text>
            </Group>
            <Button variant="subtle" onClick={handleLogout}>
              Выход
            </Button>
          </Group>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        <NavLink
          component={Link}
          to="/"
          label="Главная"
          onClick={() => toggle()}
        />
        <NavLink
          component={Link}
          to="/subjects"
          label="Предметы"
          onClick={() => toggle()}
        />
        <NavLink
          component={Link}
          to="/projects"
          label="Проекты"
          onClick={() => toggle()}
        />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
};

