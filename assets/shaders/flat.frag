#version 330 core

in vec3 pass_local_normal;
in vec3 pass_global_normal;
in vec3 pass_position;

out vec4 out_color;

uniform mat4 u_transform;

const vec3 LightColor = vec3(1.);
const vec3 LightDir = normalize(vec3(0,0.5,1));
const float AmbientLight = 1./3.;

void main()
{
    vec3 dx = dFdx(pass_position.xyz);
    vec3 dy = dFdy(pass_position.xyz);
    vec3 normal = pass_global_normal; //normalize(cross(dx, dy));
    float brightness = dot(normal, LightDir);
    brightness = clamp(brightness, 0, 1);
    // out_color.rgb = pass_local_normal*0.5+vec3(.5);
    out_color.rgb = LightColor * (brightness*(1.-AmbientLight)+AmbientLight);
    out_color.a = 1.0;
} 