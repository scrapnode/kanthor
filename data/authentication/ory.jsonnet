local session = std.extVar('session');
{
	sub: session.identity.id,
	name: if std.length(session.identity.traits.username) != 0 then session.identity.traits.username else session.identity.traits.email,
	metadata: {
		[if std.type(session.identity.metadata_public) == 'object' && 'avatar' in session.identity.metadata_public then 'avatar' else null]: session.identity.metadata_public.avatar,
	}
}