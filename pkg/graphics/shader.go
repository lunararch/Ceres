package graphics

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	program uint32
}

func NewShader(vertexSource, fragmentSource string) (*Shader, error) {
    vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
    if err != nil {
        return nil, fmt.Errorf("vertex shader compilation failed: %w", err)
    }
    defer gl.DeleteShader(vertexShader)

    fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
    if err != nil {
        return nil, fmt.Errorf("fragment shader compilation failed: %w", err)
    }
    defer gl.DeleteShader(fragmentShader)

    program := gl.CreateProgram()
    gl.AttachShader(program, vertexShader)
    gl.AttachShader(program, fragmentShader)
    gl.LinkProgram(program)

    var status int32
    gl.GetProgramiv(program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

        return nil, fmt.Errorf("shader program linking failed: %s", log)
    }

    return &Shader{program: program}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source + "\x00")
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        shaderTypeName := "unknown"
        switch shaderType {
        case gl.VERTEX_SHADER:
            shaderTypeName = "vertex"
        case gl.FRAGMENT_SHADER:
            shaderTypeName = "fragment"
        }

        return 0, fmt.Errorf("%s shader compilation error: %s", shaderTypeName, log)
    }

    return shader, nil
}

func (s *Shader) Use() {
    gl.UseProgram(s.program)
}

func (s *Shader) Delete() {
    gl.DeleteProgram(s.program)
}

func (s *Shader) GetUniformLocation(name string) int32 {
    return gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
}

func (s *Shader) SetInt(name string, value int32) {
    gl.Uniform1i(s.GetUniformLocation(name), value)
}

func (s *Shader) SetFloat(name string, value float32) {
    gl.Uniform1f(s.GetUniformLocation(name), value)
}

func (s *Shader) SetVec3(name string, x, y, z float32) {
    gl.Uniform3f(s.GetUniformLocation(name), x, y, z)
}

func (s *Shader) SetVec4(name string, x, y, z, w float32) {
    gl.Uniform4f(s.GetUniformLocation(name), x, y, z, w)
}

func (s *Shader) SetMat4(name string, mat *float32) {
    gl.UniformMatrix4fv(s.GetUniformLocation(name), 1, false, mat)
}

func (s *Shader) GetAttribLocation(name string) uint32 {
    return uint32(gl.GetAttribLocation(s.program, gl.Str(name+"\x00")))
}