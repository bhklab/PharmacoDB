class Profile < ApplicationRecord
	belongs_to :experiments
	has_many :cell, through: :experiments
	has_many :drug, through: :experiments
	has_many :dataset, through: :experiments
end
