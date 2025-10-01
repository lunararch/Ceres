package graphics

import (
	"fmt"
	"io/ioutil"
)

type ShaderManager struct {
	shaders map[string]*Shader
}

func NewShaderManager() *ShaderManager {
	return &ShaderManager{
		shaders: make(map[string]*Shader),
	}
}

func (sm *ShaderManager) LoadShader(name, vertexSource, fragmentSource string) error {
	shader, err := NewShader(vertexSource, fragmentSource)
	if err != nil {
		return fmt.Errorf("failed to load shader '%s': %w", name, err)
	}

	sm.shaders[name] = shader
	fmt.Printf("✓ Shader '%s' loaded successfully\n", name)
	return nil
}

func (sm *ShaderManager) LoadShaderFromFile(name, vertexPath, fragmentPath string) error {
	vertexSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return fmt.Errorf("failed to read vertex shader file: %w", err)
	}

	fragmentSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return fmt.Errorf("failed to read fragment shader file: %w", err)
	}

	return sm.LoadShader(name, string(vertexSource), string(fragmentSource))
}

func (sm *ShaderManager) GetShader(name string) (*Shader, error) {
	shader, exists := sm.shaders[name]
	if !exists {
		return nil, fmt.Errorf("shader '%s' not found", name)
	}
	return shader, nil
}

func (sm *ShaderManager) UseShader(name string) error {
	shader, err := sm.GetShader(name)
	if err != nil {
		return err
	}
	shader.Use()
	return nil
}

func (sm *ShaderManager) DeleteAll() {
	for name, shader := range sm.shaders {
		shader.Delete()
		fmt.Printf("✓ Shader '%s' deleted\n", name)
	}
	sm.shaders = make(map[string]*Shader)
}

func (sm *ShaderManager) LoadDefaultShaders() error {
	if err := sm.LoadShader("basic", BasicVertexShader, BasicFragmentShader); err != nil {
		return err
	}

	if err := sm.LoadShader("simple", SimpleVertexShader, SimpleFragmentShader); err != nil {
		return err
	}

	if err := sm.LoadShader("unlit", UnlitVertexShader, UnlitFragmentShader); err != nil {
		return err
	}

	fmt.Println("✓ All default shaders loaded")
	return nil
}
