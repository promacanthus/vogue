# 01-logs

## 问题

```bash
failed to create fsnotify watcher: too many open files
```

这是因为系统默认的`fs.inotify.max_user_instances=128`太小，重新设置此值：

```bash
sudo sysctl fs.inotify.max_user_instances=8192
```

## 原理
