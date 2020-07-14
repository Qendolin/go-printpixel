#version 330 core
layout (location = 0) in vec2 in_position;

out vec2 pass_uv;

uniform mat3 u_transform;

void main()
{
    gl_Position = vec4(vec3(in_position, 0.) * u_transform, 1.);
    pass_uv = in_position*vec2(1,-1)*0.5+vec2(0.5);
}