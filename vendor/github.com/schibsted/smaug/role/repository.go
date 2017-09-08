package role

import (
	"github.com/go-errors/errors"
	"github.com/go-ini/ini"
)

type RoleRepository interface {
	FindRoleByJobId(string) (string, error)
}

// InMemory Role Repository
func NewInMemoryRoleRepository() *InMemoryRoleRepository {
	return &InMemoryRoleRepository{}
}

type InMemoryRoleRepository struct {
	roles map[string]string
}

func (r *InMemoryRoleRepository) AddRole(jobId string, roleArn string) {
	if r.roles == nil {
		r.roles = make(map[string]string)
	}
	r.roles[jobId] = roleArn
}
func (r *InMemoryRoleRepository) FindRoleByJobId(jobId string) (string, error) {
	if role, ok := r.roles[jobId]; ok {
		return role, nil
	}

	return "", errors.Errorf("Role for job %s do not exist", jobId)
}

type FileRoleRepository struct {
	path  string
	roles map[string]string
}

// File Role Repository
func NewFileRoleRepository(file string) (*FileRoleRepository, error) {
	roles := make(map[string]string)
	repository := &FileRoleRepository{file, roles}
	err := repository.loadRolesFromFile()

	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *FileRoleRepository) loadRolesFromFile() error {
	loader := NewIniFileLoader(r.path)
	roles, err := loader.Load()

	if err != nil {
		return err
	}

	r.roles = roles
	return nil
}

func (r *FileRoleRepository) FindRoleByJobId(jobId string) (string, error) {
	if role, ok := r.roles[jobId]; ok {
		return role, nil
	}

	return "", errors.Errorf("Role for job %s do not exist", jobId)
}

type FileLoader interface {
	Load() (map[string]string, error)
}

// Ini File Loader
func NewIniFileLoader(file string) *IniFileLoader {
	return &IniFileLoader{file}
}

type IniFileLoader struct {
	path string
}

func (l *IniFileLoader) Load() (map[string]string, error) {
	cfg, err := ini.Load(l.path)

	if err != nil {
		return nil, err
	}

	section := cfg.Section("roles")
	keys := section.Keys()

	roles := make(map[string]string)
	for _, key := range keys {
		roles[key.Name()] = key.Value()
	}

	return roles, nil
}
