#version 330 core
layout (location = 0) in vec2 in_position;
out vec3 pass_color;

void main()
{
    gl_Position = vec4(in_position, 0., 1.);
    pass_color = vec3(in_position, 0.);
}