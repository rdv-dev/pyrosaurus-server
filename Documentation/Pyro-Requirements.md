# Requirements

# Game Simulation Requirements
One of the philosophies of this project is to maintain fairness of the rules that are implemented in the existing Game and translate that to the pyrosaurus-server .
The following sections go into detail of each of the constraints which define the Pyrosaurus game.

## Win conditions

Each team has one or more Queens (except in the Nulla Regina
Arena).  Queens can be identified by unique head markings. You win
when your team destroys any opponent Queen. (In the Nulla Regina
Arena, the last man standing wins).

Contests have a time limit. If all Queens are standing when time
runs out, the contest is declared a draw and neither team
advances.

In the first arena, the opposing Queen is not hard to find but
the later arenas are larger and finding her will be challenging.

## Dino Vision
The field of vision cone shows the area the dino actually
sees. But, even if something is in his field of vision, he
may not actually see it if he turns his head too quickly. A
wide field of vision increases his chances of seeing
something because the object will be in his field of vision
longer. Of course, a wider field of vision uses more
points.

Switching to "Prey" changes the shape of the dino's head
and puts his eyes on the side of his head. This gives him a
better chance of spotting something but he can't see
directly in front of him. This hinders his ability to aim.
To compensate, prey heads can see, hear, smell and shoot
20% farther than predators.

While we are on the subject of fire, the fire ball that is
hovering in front of the dino shows the distance that his
fire will travel.  Keep this in mind when adjusting the
dino's sight. If he can shoot fire farther than he can see,
his shots won't be very accurate because he must rely on
his hearing and smell to determine where to aim.

Your dinos will shoot at any opponent within fire range. If
he can shoot farther than he can see, your dino will have to
rely on hearing and smell to aim and he will be much less
accurate.

During a contest, your dino will
remember where another dino is for about 10 seconds after that
dino is no longer in his sight.  If he looks at that dino again
within 10 seconds, the dino will be seen instantly.


## Dino Hearing
Move the slider up and down to adjust his hearing range.
When a dino hears something, his accuracy is +/- 20 degrees
This means that when he hears something, his
hearing could be telling him that the source of the sound
is up to 20 degrees away from where it really is.

Hearing range also determines how far away he can hear
calls from other dinos on his team.  He can hear these
calls up to twice as far as his hearing range.

A dino that is not walking can only be heard at half the
hearing range.


## Dino Smell
The accuracy for smell is +/- 30 degrees. This means that
when he smells something, he knows just the general
direction of the smell.  This direction could be off by as
much as 30 degrees. A dino can smell other team members,
opponents, and food.


## Firing

Your dinos will shoot at any opponent within fire range. If
he can shoot farther than he can see, your dino will have to
rely on hearing and smell to aim and he will be much less
accurate.

Fire gradually decreases in strength over time once it leaves
the dino's mouth. When fire hits a target, fast fire will do
more damage than slow fire because it has had less time to
dissipate.

If speed is fast
enough, the fire may reach its range limit before it fully
dissipates.  This will cause the remaining fire to explode.

### Density 
Determines how many individual fire balls leave the
dino's mouth at a time.  With a low density, you will usually
see distinct fire balls. With a high density, several smaller
fire balls may be interlinked.

### Pattern 
- Adjusts the tendency of the fire to spread. Low
patterns produce narrow streams and a large pattern setting
results in a wide spread. Although your dinos are usually
accurate when they shoot, the shot could miss if the other dino
is moving quickly. Your dino has a better chance of hitting his
opponent if you give his fire some spread.

### Variation 
- The bottom half of the variation scale introduces
twists to the fire creating interesting patterns. The top half
of the scale reduces the twists while adding variations in the
fire's speed. This causes some of the fire to travel a greater
distance than normal while other parts won't travel as far as
they otherwise would.

## Dino Size 
(set when you create an individual dino) determines
how much fire the dino can breath out in one breath. A large
dino can exhale more flames than a small one.


## Predator/Prey
To compensate, prey heads can see, hear, smell and shoot 20%
farther than predators.




## Neck
If he turns his head too fast and/or
his field of vision is too narrow, he may not notice
another dino. The size of the other dino and his distance
also affect how likely he will be seen.

Variety: Controls how often your dino looks away from the
direction that he is facing.  It is important that your
dino scans around for danger but he should also pay
attention to where he is going.

The disadvantage of a long neck (besides cost) is that it
makes your dino a larger target. Necks have thinner skin
than the body. Any hits to the neck will cause twice the
damage of a body shot.



## Head Size
The size of the head affects how fast the dino exhales when
breathing fire. A small head has a small mouth which can
only expel flames in a thin stream.
The same amount of fire comes out no matter what his head
size is but the length of time that it takes to come out
varies.
Head Size limits the rate that the fire can leave his mouth.

## Heart Size
A tired dino cannot
run or shoot as often.
A dino becomes tired as he runs, jumps, and shoots. A dino
that is creeping or not moving will become rested.


tail size
The tail is used for balance. A dino can turn quicker with
a long tail.


## Legs
A dino with long legs takes big steps and consequently
walks faster than one with short legs. Large legs are also
stronger so he will be able to run faster and jump farther.



## Feet
The fight arena is a muddy bog. The type and size of feet
determine how much traction your dino has when he is
running and jumping.

A small hoof provides the least traction. The largest hoof
has slightly less traction than the smallest webbed foot.
The largest webbed foot has slightly less traction than the
smallest clawed foot.  The cost of the feet gives you an
indication of the traction they give.  The more the feet
cost, the better the traction your dino will have and the
better he can jump and run.


## Skin / Health
Snake skin is extra thin, 75% as thick as a biped's
skin.
The thickness of the skin determines how much damage a dino can
absorb before he dies. The thicker the skin, the tougher to kill.
The skin on a quadruped is extra tough (like a
rhino).  His skin is 50% thicker than a biped's
skin.
During a contest, you will see skin get thinner when it is hit by
fire. The thicker you set skin here, the longer it will last in
battle.
Snakes have 75% the skin thickness of bipeds at the same skin
setting.  Quadrupeds have 50% thicker skin then bipeds.

## Legs / Endurance

A biped is agile and fast but not
as fast as a quadruped.  Bipeds are the only dinos
that can jump. 

Quadrupeds have high strength and endurance


Sprawling legs stick out to the side like a
crocodile and are slow, clumsy and cause the owner
to tire quickly. 


## Risk
Risk - controls the two lines that form a cone in front of your
dino. If, during a contest, the center of a friendly dino is in
the cone (between the two lines) then your dino will not shoot.


## Resolve
A dino with fight training will make his fight moves
anytime that he is in fight range of an opponent. He cannot
make other decisions while he is in fight range (although he
will go for food if he needs to). The Resolve slider is a way
of telling the dino when he should stop his fight moves and
make a decision even though he is in fight range.
The Resolve slider specifies the maximum pause
between shots that will allow your dino to continue to fight.
Keep in mind that your dino will still shoot at any opponent
within range even when he does not want to fight.
Your dino will try to quit the fight
when the pause between his shots is longer than the timer
amount.
Resolve also determines how close an opponent can be to food
before your dino can attempt to go to the food. With a high
Resolve, your dino will go to food even if an oppenent is close
by. If Resolve is set to low, the opponent must be far away
before your dino will go to food.


## Movements
Self - Goals are placed on the arena relative to your dino's
current position. 
Other - Goals are placed on the arena relative to another dino's
position.
Map -  Goals are placed at fixed locations on the arena
regardless of where the dinos are.
Rotate - The goals are rotated on the arena to match the
direction the dino happens to be facing when he makes
the decision to use this movement.
Fixed - The goals are always orientated the same way you set
them down, no matter what direction the dino is facing.
Mobile - Same as Rotate except the goal moves as the dino moves
and turns. There can only be one goal with mobile.

If your dino is running and the last goal is too close to the
previous goal, he may overrun it before he can stop. If this
happens, he will turn and come back to it.  Also, he may pass
his goal if he cannot turn fast enough. In this case, he may
swing around and try to reach it again or may just go on to the
next goal.

A dino is considered at his goal when he is
anywhere inside the goal's zone.
Large zones will
cause a dino to cut corners as he goes from one zone to the next.
Large zones also allow more leeway to reach a goal when other
dinos are close by.

## Fighting
Fight Training is optional.
A dino with no
Fight Training must rely on his Movement Training.
A dino with no fight training will continue to act on his
decisions while shooting at any opponent in range.

The | is a green line that shows the current direction of your
dino.  The / is a brown line that shows where the dino is turning.

Reaction times vary from dino to dino.  Reaction times also vary
from movement to movement.  When you command your dino to run,
it may take a few seconds for him to get up to full speed.

Each arrow points in the direction that your dino was
going when the arrow was made.

When a dino uses Fight Training, he doesn't necessarily go from
the first arrow to the second and so on in sequence. Each time
your dino completes a movement, he will look for the most
appropriate movement to use next. If the other dino doesn't move,
your dino probably will proceed in the sequence that you defined.
But in combat, your opponent rarely stands still and so your dino
has to be flexible to react quickly to changing circumstances. For
example, your dino is along side his opponent and uses a learning
arrow that says to side step behind his opponent. Suddenly, his
opponent turns away from him. Now your dino finds himself behind
his opponent and so he chooses the most appropriate movement for
being behind his opponent which may be to move forward.



If your dino is using Fight Training and no Learning Arrows are
near, he will be forced to pick the closest arrow even though it
may not be appropriate.  As an example, if all of the arrows are
in front of the other dino and your dino is behind him, your dino
will use a front arrow which will make him act as if he was in
front of the other dino.  This could cause him to turn to face in
the direction that he should face if he was in front of the other
dino but because he is behind the dino, he will face away from it.

You can train your dino at distances greater than his fire
range. He will still perform the movements but he will not shoot
until he is close enough to hit an opponent.


## Decisions
All of the non-dimmed (on) conditions must be true before this
decision can be made.
You can have decisions that have no conditions. Your dino
will make one of these decisions if no other decision can be
made.


## General
As you watch a contest, pay attention to your dinos' skin. As the
markings get smaller, the skin becomes thinner.  This is a good
way of determining the condition of your dinos.
If a dino isn't moving or shooting as fast as he should, zoom in
close and watch his stomach.  If he is breathing hard and fast
then he is tired.
