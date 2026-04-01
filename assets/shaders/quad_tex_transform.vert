#version 330 core
layout (location = 0) in vec2 in_position;

out vec2 pass_uv;

uniform mat4 u_transform;

void main()
{
    gl_Position = vec4(in_position, 0., 1.)*u_transform;
    gl_Position.z = clamp(gl_Position.z, -1, 1); //clamp to prevent rounding errors like 1.00001
    pass_uv = in_position*vec2(1,-1)+vec2(0.5);
}